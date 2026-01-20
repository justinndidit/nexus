package com.justinndidit.nexus.account.dtos;

import java.util.UUID;

public record  TransferEvent(
  UUID eventId,
  TransferEventPayload payload
) {

}
