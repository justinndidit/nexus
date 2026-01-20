package com.justinndidit.nexus.account.controller;

import java.util.UUID;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import com.justinndidit.nexus.account.dtos.AccountDTO;
import com.justinndidit.nexus.account.dtos.HttpResponse;
import com.justinndidit.nexus.account.dtos.TransactionDTO;
import com.justinndidit.nexus.account.service.AccountService;

import lombok.RequiredArgsConstructor;

@RestController
@RequiredArgsConstructor
public class AccountController {
  private final AccountService accountService;

  @GetMapping("/account/{account_id}")
  public ResponseEntity<HttpResponse>getAccountById(@PathVariable(name = "account_id", required=true) UUID accountId){
    AccountDTO accountData = accountService.getAccountById(accountId);

    return ResponseEntity.ok(
      new HttpResponse(
        "success",
        "account retrieved successfully",
        accountData,
        null
      )
    );
  }

  @GetMapping("/account/transaction/{transaction_id}")
  public ResponseEntity<HttpResponse> getTransactionById(@PathVariable(name = "transaction_id", required=true) UUID transactionId){
    TransactionDTO transaction = accountService.getTransactionById(transactionId);

    return ResponseEntity.ok(
      new HttpResponse(
        "success",
        "transaction retrieved successfully",
        transaction,
        null
      )
    );
  }
}
