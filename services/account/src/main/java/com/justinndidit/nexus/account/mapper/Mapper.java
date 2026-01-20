package com.justinndidit.nexus.account.mapper;

import org.springframework.stereotype.Component;

import com.justinndidit.nexus.account.domain.Account;
import com.justinndidit.nexus.account.domain.Transaction;
import com.justinndidit.nexus.account.dtos.AccountDTO;
import com.justinndidit.nexus.account.dtos.TransactionDTO;

@Component
public class Mapper {

  public AccountDTO accountModelToDTO(Account account){
    return new AccountDTO(
      account.getId(),
      account.getUserId(),
      account.getProfileId(),
      account.getAccountNumber(),
      account.getCurrency(),
      account.getAccountType(),
      account.getAccountStatus(),
      account.getAvailableBalance()
    );
  }

  public TransactionDTO transactionModelToDTO(Transaction transaction) {
    return new TransactionDTO(
      transaction.getId(),
      transaction.getFromAccountId(),
      transaction.getDestinationAccountId(),
      transaction.getCurrencyCode(),
      transaction.getAmount(),
      transaction.getCreatedAt()
    );
  }
}
