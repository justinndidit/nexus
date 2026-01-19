package com.justinndidit.nexus.account.config.dtos;

import java.math.BigDecimal;
import java.util.UUID;

public record TransferEventPayload(
  UUID fromAccountId,
  UUID destinationAccountId,
  BigDecimal amount,
  String currency_code
) {

}
