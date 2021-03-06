-- MySQL Script generated by MySQL Workbench
-- mar. 01 déc. 2020 22:02:37
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema huetify
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema huetify
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `huetify` ;
USE `huetify` ;

-- -----------------------------------------------------
-- Table `huetify`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `huetify`.`user` (
  `user_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(45) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  UNIQUE INDEX `user_id_UNIQUE` (`user_id` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `huetify`.`role`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `huetify`.`role` (
  `role_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`role_id`),
  UNIQUE INDEX `role_id_UNIQUE` (`role_id` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `huetify`.`user_role`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `huetify`.`user_role` (
  `user_id` INT UNSIGNED NOT NULL,
  `role_id` INT UNSIGNED NOT NULL,
  INDEX `fk_user_role_user_idx` (`user_id` ASC) VISIBLE,
  INDEX `fk_user_role_role1_idx` (`role_id` ASC) VISIBLE,
  CONSTRAINT `fk_user_role_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `huetify`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_role_role1`
    FOREIGN KEY (`role_id`)
    REFERENCES `huetify`.`role` (`role_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `huetify`.`item_type`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `huetify`.`item_type` (
  `type_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`type_id`),
  UNIQUE INDEX `type_id_UNIQUE` (`type_id` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `huetify`.`item`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `huetify`.`item` (
  `item_uid` VARCHAR(255) NOT NULL,
  `type_id` INT UNSIGNED NOT NULL,
  `icon` VARCHAR(255) NOT NULL DEFAULT 'default.svg',
  UNIQUE INDEX `item_id_UNIQUE` (`item_uid` ASC) VISIBLE,
  INDEX `fk_item_item_type1_idx` (`type_id` ASC) VISIBLE,
  CONSTRAINT `fk_item_item_type1`
    FOREIGN KEY (`type_id`)
    REFERENCES `huetify`.`item_type` (`type_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `huetify`.`sensor_config`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `huetify`.`sensor_config` (
  `item_uid` INT NOT NULL,
  `refresh_rate` INT NOT NULL DEFAULT 2000,
  UNIQUE INDEX `item_uid_UNIQUE` (`item_uid` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `huetify`.`group`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `huetify`.`group` (
  `group_hid` VARCHAR(255) NOT NULL,
  `position` INT NOT NULL DEFAULT 1,
  UNIQUE INDEX `group_hid_UNIQUE` (`group_hid` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `huetify`.`app`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `huetify`.`app` (
  `name` VARCHAR(255) NOT NULL DEFAULT 'Huetify',
  `language` VARCHAR(45) NOT NULL DEFAULT 'en',
  `analytics` TINYINT NOT NULL DEFAULT 1)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
