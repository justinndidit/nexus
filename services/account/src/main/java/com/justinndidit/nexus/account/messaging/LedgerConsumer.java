package com.justinndidit.nexus.account.messaging;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import com.justinndidit.nexus.account.config.CustomLogger;
import com.justinndidit.nexus.account.domain.Account;
import com.justinndidit.nexus.account.domain.ProcessedEvent;
import com.justinndidit.nexus.account.domain.Transaction;
import com.justinndidit.nexus.account.dtos.TransferEvent;
import com.justinndidit.nexus.account.repository.AccountRepository;
import com.justinndidit.nexus.account.repository.ProcessedEventRepository;
import com.justinndidit.nexus.account.repository.TransactionRepository;

import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;

@Component
@RequiredArgsConstructor
public class LedgerConsumer {
// private final ObjectMapper objectMapper;

  private final ProcessedEventRepository eventRepository;
  private final AccountRepository accountRepository;
  private final TransactionRepository transactionRepository;
  private final CustomLogger logger;

/*
  @KafkaLister designates a bean methof as a message listener

*/

  @KafkaListener(topics="ledger.transactions.v1", groupId="com.justinndidit.nexus.transactions")
  @Transactional
  public void consume(TransferEvent message) {
    try {

      if (!eventRepository.findById(message.eventId()).isEmpty()) {
        logger.warnWithArguments("event {} already processed", message.eventId());
        return;
      }
      updateBalance(message.payload().fromAccountId(), message.payload().amount().negate());
      updateBalance(message.payload().destinationAccountId(), message.payload().amount());
      eventRepository.save(new ProcessedEvent(message.eventId(),LocalDateTime.now()));

      transactionRepository.save(new Transaction(
        message.payload().fromAccountId(),
        message.payload().destinationAccountId(),
        message.payload().currency_code(),
        message.payload().amount(),
        LocalDateTime.now()
      ));

      logger.infoWithArguments("event {} processed successfully", message.eventId());
    } catch(Exception e) {
      logger.errorWithArguments("{}: {}",e.getClass(),e.getMessage());
      throw new RuntimeException(e);
    }
  }

  public void updateBalance(UUID accountId, BigDecimal amount) {
    Account account = accountRepository.
                        findById(accountId).
                          orElseThrow(() -> new RuntimeException("Account not found: " + accountId));
    account.setAvailableBalance(account.getAvailableBalance().add(amount));
    accountRepository.save(account);
  }
}