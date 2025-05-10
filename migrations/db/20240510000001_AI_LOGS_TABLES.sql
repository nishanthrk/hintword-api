-- +goose Up

--
-- Table structure for table `notes`
--
CREATE TABLE note_ai_logs (
    `log_id` CHAR(36) COLLATE utf8mb4_general_ci NOT NULL,
    `note_id` CHAR(36) COLLATE utf8mb4_general_ci NOT NULL,
    `user_id` CHAR(36) COLLATE utf8mb4_general_ci NOT NULL,
    `interaction_type` ENUM('COMPLETION', 'STT', 'TTS') NOT NULL,
    `request_payload` LONGTEXT COLLATE utf8mb4_general_ci DEFAULT NULL,      -- JSON input: prompt, audio blob info, etc.
    `response_payload` LONGTEXT COLLATE utf8mb4_general_ci DEFAULT NULL,     -- JSON response: completion, transcript, etc.
    `token_usage` JSON DEFAULT NULL, -- {"prompt_tokens": 100, "completion_tokens": 150, "total_tokens": 250}
    `audio_url` VARCHAR(512) COLLATE utf8mb4_general_ci DEFAULT NULL,  -- if audio was involved (TTS or STT)
    `model_used` VARCHAR(100) COLLATE utf8mb4_general_ci DEFAULT NULL, -- e.g., 'gpt-4', 'whisper', 'elevenlabs-tts'
    `language_code` VARCHAR(10) COLLATE utf8mb4_general_ci DEFAULT NULL, -- e.g., 'en-US'
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`log_id`),
    CONSTRAINT `fk_note_ai_logs_notes` FOREIGN KEY (`note_id`) REFERENCES `notes`(`note_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_note_ai_logs_users` FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- +goose Down

DROP TABLE IF EXISTS `note_ai_logs`;