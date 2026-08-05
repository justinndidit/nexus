package com.justinndidit.nexus.account;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import io.github.cdimascio.dotenv.Dotenv;

@SpringBootApplication
public class AccountApplication {
  public static void main(String[] args) {
    Dotenv dotenv = Dotenv.configure().load();
    dotenv.entries().forEach(e -> System.setProperty(e.getKey(), e.getValue()));

    SpringApplication.run(AccountApplication.class, args);
  }
}
