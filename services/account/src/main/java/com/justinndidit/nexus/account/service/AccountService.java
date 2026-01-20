package com.justinndidit.nexus.account.service;

import java.util.UUID;

import com.justinndidit.nexus.account.dtos.AccountDTO;

public interface AccountService {
  public AccountDTO getAccountById(UUID accountId);
}
