package com.justinndidit.nexus.account.messaging;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import com.justinndidit.nexus.account.config.dtos.TransferEvent;
import com.justinndidit.nexus.account.domain.Account;
import com.justinndidit.nexus.account.domain.ProcessedEvent;
import com.justinndidit.nexus.account.repository.AccountRepository;
import com.justinndidit.nexus.account.repository.ProcessedEventRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import tools.jackson.core.JacksonException;
// import tools.jackson.databind.ObjectMapper;

@Component
@RequiredArgsConstructor
public class LedgerConsumer {
// private final ObjectMapper objectMapper;

  private final ProcessedEventRepository eventRepository;
  private final AccountRepository accountRepository;

/*
  @KafkaLister designates a bean methof as a message listener

*/

  @KafkaListener(topics="ledger.transactions.v1", groupId="com.justinndidit.nexus.transactions")
  @Transactional
  public void consume(TransferEvent message) {
    try {

      if (!eventRepository.findById(message.eventId()).isEmpty()) {
        //log already processed event
        //return
      }

      updateBalance(message.payload().fromAccountId(), message.payload().amount().negate());
      updateBalance(message.payload().destinationAccountId(), message.payload().amount());
      eventRepository.save(new ProcessedEvent(message.eventId(),LocalDateTime.now() ));

    } catch (JacksonException e) {
      System.out.println(e.toString());
    }catch(Exception e) {
      System.out.println(e.toString());
    }
  }

  public void updateBalance(UUID accountId, BigDecimal amount) {
    Account account = accountRepository.findById(accountId).orElseThrow(() -> new RuntimeException("Account not found: " + accountId));
    account.setAvailableBalance(amount);
    accountRepository.save(account);
  }
}