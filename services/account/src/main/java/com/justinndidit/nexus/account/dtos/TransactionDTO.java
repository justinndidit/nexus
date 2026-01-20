package com.justinndidit.nexus.account.dtos;

import java.util.UUID;
import java.math.BigDecimal;
import java.time.LocalDateTime;

public record TransactionDTO(
  UUID id,
  UUID fromAccountId,
  UUID destinationAccountId,
  String currency,
  BigDecimal amount,
  LocalDateTime createdAt
) {}
