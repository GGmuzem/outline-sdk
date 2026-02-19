import asyncio
import logging
import uuid
from dataclasses import dataclass
from enum import Enum
from typing import Optional

from sqlalchemy.ext.asyncio import AsyncSession
from yookassa import Configuration, Payment

from config import settings
from models import SubscriptionTier, User

logger = logging.getLogger(__name__)

# Настройка YooKassa
Configuration.account_id = settings.yookassa_shop_id
Configuration.secret_key = settings.yookassa_secret_key


class PaymentStatus(str, Enum):
    PENDING = "pending"
    WAITING_FOR_CAPTURE = "waiting_for_capture"
    SUCCEEDED = "succeeded"
    CANCELED = "canceled"

    def is_successful(self) -> bool:
        return self == PaymentStatus.SUCCEEDED


@dataclass
class PaymentData:
    yookassa_payment_id: str
    status: PaymentStatus
    subscription_tier: SubscriptionTier
    confirmation_url: Optional[str] = None
    amount: float = 0.0


class YooKassaService:
    async def create_payment(self, session: AsyncSession, user: User, tier: SubscriptionTier) -> PaymentData:
        """Создание платежа в ЮKassa."""
        price = settings.get_tier_price(tier.value)
        idempotence_key = str(uuid.uuid4())
        description = f"Подписка {tier.value} для пользователя {user.id}"

        def _create():
            return Payment.create({
                "amount": {
                    "value": str(price),
                    "currency": "RUB"
                },
                "confirmation": {
                    "type": "redirect",
                    "return_url": settings.yookassa_return_url
                },
                "capture": True,
                "description": description,
                "metadata": {
                    "user_id": user.id,
                    "tier": tier.value
                }
            }, idempotence_key)

        # Выполняем синхронный запрос в отдельном потоке
        payment_response = await asyncio.to_thread(_create)

        return PaymentData(
            yookassa_payment_id=payment_response.id,
            status=PaymentStatus(payment_response.status),
            subscription_tier=tier,
            confirmation_url=payment_response.confirmation.confirmation_url,
            amount=float(payment_response.amount.value)
        )

    async def check_payment(self, session: AsyncSession, payment_id: str) -> PaymentData:
        """Проверка статуса платежа."""
        def _find():
            return Payment.find_one(payment_id)

        payment_response = await asyncio.to_thread(_find)
        
        tier_str = payment_response.metadata.get("tier", "free")
        try:
            tier = SubscriptionTier(tier_str)
        except ValueError:
            tier = SubscriptionTier.FREE

        return PaymentData(
            yookassa_payment_id=payment_response.id,
            status=PaymentStatus(payment_response.status),
            subscription_tier=tier,
            amount=float(payment_response.amount.value)
        )