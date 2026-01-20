package com.justinndidit.nexus.account.config;

import org.springframework.stereotype.Component;
import org.tinylog.Logger;

@Component
public class CustomLogger {
  String serviceName;

  public CustomLogger(){
    this.serviceName = "account-service";
  }

  public void warn(String msg) {
    Logger.warn("{}: {}", this.serviceName, msg);
  }

  public void warnWithArguments(String msg, Object ...args) {
    String message = this.serviceName + ": " + msg;
    Logger.warn(message, args);
  }

  public void infoWithArguments(String msg, Object ...args) {
    String message = this.serviceName + ": " + msg;
    Logger.info(message, args);
  }

  public void info(String msg) {
    Logger.info("{}: {}", this.serviceName, msg);
  }

  public void error(String msg) {
    Logger.error("{}: {}", this.serviceName, msg);
  }

  public void errorWithArguments(String msg, Object ...args) {
    String message = this.serviceName + ": " + msg;
    Logger.error(message, args);
  }

  public void errorWithException(Exception e,String msg) {
    String message = this.serviceName + ": " + msg;
    Logger.error(e, message);
  }


}
