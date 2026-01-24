import logging

class CustomLogger:
  def __init__(self):
    self.logger = logging.getLogger("notification_service")

  def info(self,msg : str) -> None:
    self.logger.info(msg)

  def warn(self, msg : str) -> None:
    self.logger.warning(msg)

  def error(self, msg : str) -> None:
    self.logger.error(msg)

  def exception(self, msg : str, exception : Exception) -> None:
    self.logger.exception(
      msg,
      exception
    )

def get_custom_logger() -> CustomLogger:
  return CustomLogger()