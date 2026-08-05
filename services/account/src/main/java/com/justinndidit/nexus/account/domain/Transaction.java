package com.justinndidit.nexus.account.domain;

import java.math.BigDecimal;
import java.time.Instant;
import java.util.UUID;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Entity
@Table(name = "transactions")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Transaction{
  @Id
  private UUID id;

  private UUID fromAccountId;
  private UUID destinationAccountId;

  private String currencyCode;
  private BigDecimal amount;

  private Instant createdAt;

}
