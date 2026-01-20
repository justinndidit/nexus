package com.justinndidit.nexus.account.repository;

import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;

import com.justinndidit.nexus.account.domain.Transaction;

public interface TransactionRepository extends JpaRepository<Transaction, UUID>{

}
