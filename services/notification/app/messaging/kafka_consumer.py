from aiokafka import AIOKafkaConsumer, ConsumerRecord
from services.notification.app.core.logger import CustomLogger, get_custom_logger
from services.notification.app.services.email_service import EmailService, get_email_service
from services.notification.app.services.push_service import   PushService, get_push_service
from services.notification.app.services.sms_service import SMS_Service,  get_sms_service

from fastapi import Depends

class KafkaConsumer:
  def __init__(self,
               kafka_address : str,
               email_service : EmailService,
               sms_service : SMS_Service,
               push_service : PushService,
               logger : CustomLogger ):
    self.logger = logger
    self.email_service = email_service
    self.sms_service = sms_service
    self.push_service = push_service
    self.kafka_address = kafka_address
    self.base_group_id = "justinndidit.nexus.notification"

  async def _consume_notification(self, topic : str, group_id : str) -> AIOKafkaConsumer | None:
    return AIOKafkaConsumer(
      topic,
      bootstrap_servers=self.kafka_address,
      group_id = group_id
    )

  async def consume_ledger_transaction(self) -> None:
    consumer : AIOKafkaConsumer = await self._consume_notification(
      topic="ledger.transactions.v1",
      group_id=self.base_group_id + "ledger_transactions"
    )
    await consumer.start()
    try:
      async for msg in consumer:
        self.logger.info(f"message {msg.key} received")
        pass
    except Exception as e:
      self.logger.error(f"error consuming kafka message: {str(e)}")
      return None
    finally:
      consumer.stop()