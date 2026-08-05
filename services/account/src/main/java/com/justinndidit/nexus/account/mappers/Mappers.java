package com.justinndidit.nexus.account.mappers;

import org.springframework.stereotype.Component;

import com.justinndidit.nexus.account.domain.Account;
import com.justinndidit.nexus.account.domain.Transaction;
import com.justinndidit.nexus.account.dtos.AccountDTO;
import com.justinndidit.nexus.account.dtos.TransactionDTO;

@Component
public class Mappers {

  public AccountDTO accountModelToDTO(Account account){
    return new AccountDTO(
      account.getId(),
      account.getUserId(),
      account.getProfileId(),
      account.getAccountNumber(),
      account.getCurrencyCode(),
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
