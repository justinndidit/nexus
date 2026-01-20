package com.justinndidit.nexus.account.domain;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.Version;
import lombok.Data;

@Entity
@Table(name = "account")
@Data
public class Account {
  @Id
  @GeneratedValue(strategy=GenerationType.UUID)
  private UUID id;

  private UUID userId;
  private UUID profileId;

  private String accountNumber;
  private String currency;
  private String accountType;
  private String accountStatus;
  private BigDecimal availableBalance;
  // private BigDecimal ledgerBalance;

  @Version
  private long version; //optimistic locking

  private LocalDateTime createdAt;
  private LocalDateTime updatedAt;

}
