package com.justinndidit.nexus.account.domain;

import java.math.BigDecimal;
import java.time.Instant;
import java.util.UUID;

import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.Version;
import lombok.Data;
import lombok.NoArgsConstructor;

@Entity
@Table(name = "accounts")
@NoArgsConstructor
@Data
public class Account {
  @Id
  @GeneratedValue(strategy=GenerationType.UUID)
  private UUID id;

  private UUID userId;
  private UUID profileId;

  private String accountNumber;
  private String currencyCode;
  private String accountType;
  private String accountStatus;
  private BigDecimal availableBalance;
  // private BigDecimal ledgerBalance;

  @Version
  private long version; //optimistic locking

  @CreationTimestamp
  private Instant createdAt;

  @UpdateTimestamp
  private Instant updatedAt;

}
