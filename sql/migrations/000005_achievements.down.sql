DROP TABLE IF EXISTS achievements_completed;
DROP TABLE IF EXISTS achievements;

DROP TRIGGER IF EXISTS trg_award_achievement_xp ON achievements_completed;
DROP TRIGGER IF EXISTS trg_check_xp_achievement ON users;
DROP TRIGGER IF EXISTS trg_check_first_lesson_achievement ON lessons_completed;

DROP FUNCTION IF EXISTS fn_award_achievement_xp();
DROP FUNCTION IF EXISTS fn_check_xp_achievement();
DROP FUNCTION IF EXISTS fn_check_first_lesson_achievement();

