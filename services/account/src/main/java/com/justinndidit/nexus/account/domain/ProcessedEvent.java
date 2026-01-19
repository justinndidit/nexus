package com.justinndidit.nexus.account.domain;

import java.time.LocalDateTime;
import java.util.UUID;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Data;

@Entity
@Table(name = "processed_event")
@Data
@AllArgsConstructor
public class ProcessedEvent {
  @Id
  private UUID Id;
  private LocalDateTime createdAt;
}
