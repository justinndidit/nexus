package com.justinndidit.nexus.account.service;

import java.util.UUID;

import com.justinndidit.nexus.account.dtos.AccountDTO;
import com.justinndidit.nexus.account.dtos.TransactionDTO;

public interface AccountService {
  public AccountDTO getAccountById(UUID accountId);
  public TransactionDTO getTransactionById(UUID transactionId);
}
