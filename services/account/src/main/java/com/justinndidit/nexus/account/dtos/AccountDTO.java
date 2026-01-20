package com.justinndidit.nexus.account.dtos;

import java.math.BigDecimal;
import java.util.UUID;

public record AccountDTO(
  UUID id,
  UUID userId,
  UUID profileId,
  String accountNumber,
  String currency,
  String accountType,
  String accountStatus,
  BigDecimal availabelBalance
) {}
