CREATE TABLE achievements (
	id INT PRIMARY KEY,
	title VARCHAR(64) NOT NULL,
	description TEXT NOT NULL,
	icon_data TEXT NOT NULL DEFAULT '',
	xp INT DEFAULT 0
);

CREATE TABLE achievements_completed (
	user_id INT NOT NULL,
	achievement_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (achievement_id) REFERENCES achievements (id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, achievement_id)
);

INSERT INTO achievements(id, title, description, xp, icon_data)
			VALUES (1, 'Первый урок',
					'Завершить свой первый урок по химии',
					120,
				'https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2F7fon.club%2Fpics%2Fuploads%2Fposts%2F2023-09%2F1694937699_7fon-club-p-vafelnaya-kartinka-luntik-10.jpg');
INSERT INTO achievements(id,title, description, xp, icon_data)
			VALUES (2, 'Старательный ученик',
					'Заработать 1000 очков опыта',
					220,
				'https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fzefirka.club%2Fwallpapers%2Fuploads%2Fposts%2F2023-07%2F1688455278_zefirka-club-p-malish-barboskini-na-belom-fone-43.jpg');
INSERT INTO achievements(id,title, description, xp, icon_data)
			VALUES (3, 'Целеустремленный химик',
					'10 дней подряд изучать химию',
					220,
				'https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Ftse3.mm.bing.net%2Fth%2Fid%2FOIP.imuB4KvoHBy7fH-EczvGlgHaHa%3Fr%3D0%26pid%3DApi&f=1&ipt=76cd4ffebaf4554bc906b086da0374a721aac4e61e94367c1c88e22f37254316');

CREATE OR REPLACE FUNCTION fn_award_achievement_xp()
RETURNS TRIGGER AS $$
DECLARE
    v_achievement_xp INT;
BEGIN
    SELECT xp INTO v_achievement_xp FROM achievements WHERE id = NEW.achievement_id;
    
    IF v_achievement_xp IS NOT NULL THEN
        UPDATE users 
        SET xp = xp + v_achievement_xp 
        WHERE id = NEW.user_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_award_achievement_xp
AFTER INSERT ON achievements_completed
FOR EACH ROW
EXECUTE FUNCTION fn_award_achievement_xp();

-- 2. Auto-grant "Старательный ученик" (ID 2) when XP reaches 1000
CREATE OR REPLACE FUNCTION fn_check_xp_achievement()
RETURNS TRIGGER AS $$
BEGIN
    -- Only fire when crossing the 1000 XP threshold
    IF OLD.xp < 1000 AND NEW.xp >= 1000 THEN
        INSERT INTO achievements_completed (user_id, achievement_id) 
        VALUES (NEW.id, 2)
        ON CONFLICT (user_id, achievement_id) DO NOTHING;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_xp_achievement
AFTER UPDATE OF xp ON users
FOR EACH ROW
EXECUTE FUNCTION fn_check_xp_achievement();

-- 3. Auto-grant "Первый урок" (ID 1) on first lesson completion
CREATE OR REPLACE FUNCTION fn_check_first_lesson_achievement()
RETURNS TRIGGER AS $$
DECLARE
    v_lesson_count INT;
BEGIN
    SELECT COUNT(*) INTO v_lesson_count 
    FROM lessons_completed 
    WHERE user_id = NEW.user_id;
    
    IF v_lesson_count = 1 THEN
        INSERT INTO achievements_completed (user_id, achievement_id) 
        VALUES (NEW.user_id, 1)
        ON CONFLICT (user_id, achievement_id) DO NOTHING;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_first_lesson_achievement
AFTER INSERT ON lessons_completed
FOR EACH ROW
EXECUTE FUNCTION fn_check_first_lesson_achievement();
