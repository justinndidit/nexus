package com.justinndidit.nexus.account.mapper;

import org.springframework.stereotype.Component;

import com.justinndidit.nexus.account.domain.Account;
import com.justinndidit.nexus.account.dtos.AccountDTO;

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

}
