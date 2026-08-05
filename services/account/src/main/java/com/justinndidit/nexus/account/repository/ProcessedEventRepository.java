package com.justinndidit.nexus.account.repository;

import java.util.UUID;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.justinndidit.nexus.account.domain.ProcessedEvents;

@Repository
public interface  ProcessedEventRepository extends JpaRepository<ProcessedEvents, UUID>{
}
