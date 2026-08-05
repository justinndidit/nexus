package com.justinndidit.nexus.account.dtos;

import java.util.UUID;
import java.math.BigDecimal;
import java.time.Instant;

public record TransactionDTO(
  UUID id,
  UUID fromAccountId,
  UUID destinationAccountId,
  String currencyCode,
  BigDecimal amount,
  Instant createdAt
) {}
