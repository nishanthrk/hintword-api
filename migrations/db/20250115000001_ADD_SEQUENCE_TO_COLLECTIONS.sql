-- +goose Up
--
-- Add sequence column to collections table for manual ordering
--

ALTER TABLE `collections` 
ADD COLUMN `sequence` BIGINT NOT NULL DEFAULT 0 AFTER `status`;

-- Initialize sequence values for existing collections based on creation time
UPDATE `collections` 
SET `sequence` = UNIX_TIMESTAMP(`created_at`) 
WHERE `sequence` = 0;

-- Add index for efficient ordering by user and sequence
CREATE INDEX `idx_collections_user_sequence` ON `collections` (`user_id`, `sequence`);

-- +goose Down

-- Remove the index
DROP INDEX `idx_collections_user_sequence` ON `collections`;

-- Remove the sequence column
ALTER TABLE `collections` 
DROP COLUMN `sequence`;
