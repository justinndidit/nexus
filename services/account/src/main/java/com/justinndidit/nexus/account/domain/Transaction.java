package com.justinndidit.nexus.account.domain;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import lombok.Data;

@Entity
@Data
public class Transaction{
  @Id
  @GeneratedValue(strategy=GenerationType.UUID)
  private UUID id;

  private UUID fromAccountId;
  private UUID destinationAccountId;

  private String currencyCode;
  private BigDecimal amount;

  private LocalDateTime createdAt;

  public Transaction(
        UUID fromAccountId,
        UUID destinationAccountId,
        String currency,
        BigDecimal amount,
        LocalDateTime created){

    this.fromAccountId = fromAccountId;
    this.destinationAccountId = destinationAccountId;
    this.currencyCode = currency;
    this.amount = amount;
    this.createdAt = created;
  }
}