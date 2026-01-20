package com.justinndidit.nexus.account.dtos;

public record  HttpResponse(
  String status,
  String message,
  Object data,
  Object error
){}
