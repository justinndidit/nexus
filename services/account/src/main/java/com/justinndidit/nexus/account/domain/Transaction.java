package com.justinndidit.nexus.account.domain;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import lombok.AllArgsConstructor;
import lombok.Data;

@Entity
@Data
@AllArgsConstructor
public class Transaction{
  @Id
  private UUID id;

  private UUID fromAccountId;
  private UUID destinationAccountId;

  private String currencyCode;
  private BigDecimal amount;

  private LocalDateTime createdAt;

}