package com.justinndidit.nexus.account.config.dtos;

import java.util.UUID;

public record  TransferEvent(
  UUID eventId,
  TransferEventPayload payload
) {

}
