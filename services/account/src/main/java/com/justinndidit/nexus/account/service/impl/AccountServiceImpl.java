package com.justinndidit.nexus.account.service.impl;

import java.util.UUID;

import org.springframework.stereotype.Service;

import java.util.Optional;


import com.justinndidit.nexus.account.config.CustomLogger;
import com.justinndidit.nexus.account.domain.Account;
import com.justinndidit.nexus.account.dtos.AccountDTO;
import com.justinndidit.nexus.account.repository.AccountRepository;
import com.justinndidit.nexus.account.service.AccountService;
import com.justinndidit.nexus.account.mapper.Mapper;

import jakarta.persistence.EntityNotFoundException;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class AccountServiceImpl implements AccountService {
  private final AccountRepository accountRepo;
  private final CustomLogger logger;
  private final Mapper mapper;

  @Override
  public AccountDTO getAccountById(UUID accountId) {

    Optional<Account> account = accountRepo.findById(accountId);
    if (account.isEmpty()){
      logger.errorWithArguments("account {} does not exist", accountId);
      throw new EntityNotFoundException("No account with id "+ accountId);
    }

    return mapper.accountModelToDTO(account.get());
  }

}
