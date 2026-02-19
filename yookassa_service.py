import asyncio
import hashlib
import hmac
import json
import logging
from datetime import datetime, timezone
from typing import Optional
from uuid import uuid4

from sqlalchemy import select, update
from sqlalchemy.ext.asyncio import AsyncSession
from yookassa import Configuration, Payment as YooPayment

from config import settings
from models import Payment, PaymentStatus, SubscriptionTier, User
from services import UserService

logger = logging.getLogger(__name__)


def _map_status(status: str) -> PaymentStatus:
    try:
        return PaymentStatus(status)
    except ValueError:
        if status == "waiting_for_capture":
            return PaymentStatus.WAITING_FOR_CAPTURE
        return PaymentStatus.PENDING


class YooKassaService:
    """YooKassa payment integration with webhook signature verification."""

    def __init__(self) -> None:
        logger.info("Initializing YooKassaService (v3.2 - DEBUG)")
        Configuration.account_id = settings.yookassa_shop_id
        Configuration.secret_key = settings.yookassa_secret_key

    def _get_tier_price_rubles(self, tier: SubscriptionTier) -> int:
        """Get price in RUB for subscription tier."""
        if tier == SubscriptionTier.PRO:
            return settings.pro_monthly_price
        if tier == SubscriptionTier.ULTRA:
            return settings.ultra_monthly_price
        return 0

    def verify_webhook_signature(self, payload: dict, signature: str) -> bool:
        """Verify YooKassa webhook signature for security.
        
        Args:
            payload: Webhook payload as dict
            signature: Signature from X-Signature header
            
        Returns:
            True if signature is valid, False otherwise
        """
        if not settings.yookassa_webhook_secret:
            logger.warning("YooKassa webhook secret not configured, skipping verification")
            return True  # Allow if not configured (development mode)

        try:
            # YooKassa uses HMAC-SHA256 for webhook signatures
            # Format: hmac_sha256(webhook_secret, notification_body)
            notification_body = json.dumps(payload, separators=(',', ':'), ensure_ascii=False)
            
            expected_signature = hmac.new(
                settings.yookassa_webhook_secret.encode('utf-8'),
                notification_body.encode('utf-8'),
                hashlib.sha256
            ).hexdigest()

            is_valid = hmac.compare_digest(expected_signature, signature)
            
            if not is_valid:
                logger.warning(
                    "Invalid webhook signature",
                    extra={"expected": expected_signature[:16], "received": signature[:16]}
                )
            
            return is_valid

        except Exception as exc:
            logger.error(f"Webhook signature verification failed: {exc}")
            return False

    async def create_payment(
        self,
        session: AsyncSession,
        user: User,
        tier: SubscriptionTier,
    ) -> Payment:
        amount_rubles = self._get_tier_price_rubles(tier)
        description = f"Подписка {tier.value.capitalize()} для пользователя {user.id}"

        if amount_rubles <= 0:
            raise ValueError(f"Invalid amount for tier {tier.value}")

        payload = {
            "amount": {"value": f"{amount_rubles:.2f}", "currency": "RUB"},
            "confirmation": {
                "type": "redirect",
                "return_url": settings.yookassa_return_url,
            },
            "capture": True,
            "description": description,
            "metadata": {
                "user_id": str(user.id),
                "tier": tier.value,
            },
        }

        try:
            idempotence_key = str(uuid4())
            response = await asyncio.to_thread(YooPayment.create, payload, idempotence_key)
        except Exception:
            logger.exception("Failed to create YooKassa payment")
            raise

        # Store amount in kopecks for consistency with payment records
        payment = Payment(
            user_id=user.id,
            yookassa_payment_id=response.id,
            subscription_tier=tier.value,
            amount=amount_rubles * 100,  # Convert to kopecks for storage
            status=_map_status(response.status),
            description=description,
            confirmation_url=response.confirmation.confirmation_url,
        )

        session.add(payment)
        await session.flush()
        return payment

    async def check_payment(
        self,
        session: AsyncSession,
        payment_id: str,
    ) -> Payment:
        result = await session.execute(
            select(Payment).where(Payment.yookassa_payment_id == payment_id)
        )
        payment: Optional[Payment] = result.scalar_one_or_none()
        if not payment:
            raise ValueError("Payment not found")

        # Если платеж уже был успешен, не проверяем API и не начисляем подписку повторно
        if payment.status == PaymentStatus.SUCCEEDED:
            return payment

        try:
            response = await asyncio.to_thread(YooPayment.find_one, payment_id)
        except Exception:
            logger.exception("Failed to check YooKassa payment")
            raise

        new_status = _map_status(response.status)
        
        # Начисляем подписку ТОЛЬКО если статус изменился на SUCCEEDED
        if new_status == PaymentStatus.SUCCEEDED and payment.status != PaymentStatus.SUCCEEDED:
            user_service = UserService(session)
            result = await session.execute(select(User).where(User.id == payment.user_id))
            user = result.scalar_one_or_none()
            if user:
                await user_service.upgrade_subscription(user, SubscriptionTier(payment.subscription_tier), duration_days=30)

        payment.status = new_status
        payment.processed_at = datetime.now(timezone.utc)
        payment.error_message = None
        if getattr(response, "cancellation_details", None):
            payment.error_message = getattr(response.cancellation_details, "reason", None)

        await session.flush()
        return payment

    async def handle_webhook(self, session: AsyncSession, payload: dict, signature: Optional[str] = None) -> Optional[Payment]:
        """Handle YooKassa webhook with signature verification.
        
        Args:
            session: Database session
            payload: Webhook payload
            signature: Optional signature from X-Signature header
            
        Returns:
            Updated Payment object or None
        """
        logger.info("Processing YooKassa webhook (v3.3 - Fixed)")
        # Verify signature if provided
        if signature and not self.verify_webhook_signature(payload, signature):
            logger.error("Webhook signature verification failed, ignoring webhook")
            return None

        event_object = payload.get("object") or {}
        payment_id = event_object.get("id")
        if not payment_id:
            return None

        # Securely verify payment status via API instead of trusting webhook payload
        try:
            return await self.check_payment(session, payment_id)
        except Exception as e:
            logger.error(f"Error handling webhook for payment {payment_id}: {e}")
            return None


__all__ = ["YooKassaService"]
