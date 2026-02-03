import sqlite3
from typing import Optional


class Database:
    def __init__(self, db_path: str = "valentines.db"):
        self.db_path = db_path
        self.init_db()

    def get_connection(self) -> sqlite3.Connection:
        """Get database connection."""
        conn = sqlite3.connect(self.db_path)
        conn.row_factory = sqlite3.Row
        return conn

    def init_db(self) -> None:
        """Initialize database schema."""
        with self.get_connection() as conn:
            conn.execute("""
                CREATE TABLE IF NOT EXISTS users (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    telegram_id INTEGER UNIQUE NOT NULL,
                    username TEXT,
                    first_name TEXT,
                    last_name TEXT,
                    is_bot INTEGER NOT NULL DEFAULT 0,
                    language_code TEXT,
                    is_premium INTEGER,
                    sex TEXT,
                    about TEXT,
                    state TEXT NOT NULL DEFAULT 'start',
                    time_ranges TEXT NOT NULL DEFAULT '000000'
                )
            """)
            conn.commit()

    def save_user(
        self,
        telegram_id: int,
        username: Optional[str],
        first_name: Optional[str],
        last_name: Optional[str],
        is_bot: bool,
        language_code: Optional[str],
        is_premium: Optional[bool],
    ) -> None:
        """Save or update user information."""
        with self.get_connection() as conn:
            conn.execute(
                """
                INSERT INTO users (
                    telegram_id, username, first_name, last_name,
                    is_bot, language_code, is_premium
                ) VALUES (?, ?, ?, ?, ?, ?, ?)
                ON CONFLICT(telegram_id) DO UPDATE SET
                    username = excluded.username,
                    first_name = excluded.first_name,
                    last_name = excluded.last_name,
                    is_bot = excluded.is_bot,
                    language_code = excluded.language_code,
                    is_premium = excluded.is_premium
                """,
                (
                    telegram_id,
                    username,
                    first_name,
                    last_name,
                    1 if is_bot else 0,
                    language_code,
                    1 if is_premium is True else (0 if is_premium is False else None),
                ),
            )
            conn.commit()

    def get_time_ranges(self, telegram_id: int) -> str:
        """Get user's time ranges as binary string (e.g., '101010')."""
        with self.get_connection() as conn:
            cursor = conn.execute(
                "SELECT time_ranges FROM users WHERE telegram_id = ?",
                (telegram_id,),
            )
            row = cursor.fetchone()
            return row["time_ranges"] if row else "000000"

    def save_time_ranges(self, telegram_id: int, time_ranges: str) -> None:
        """Save user's time ranges as binary string."""
        with self.get_connection() as conn:
            conn.execute(
                "UPDATE users SET time_ranges = ? WHERE telegram_id = ?",
                (time_ranges, telegram_id),
            )
            conn.commit()

    def get_user(self, telegram_id: int) -> Optional[sqlite3.Row]:
        """Get user by telegram_id."""
        with self.get_connection() as conn:
            cursor = conn.execute(
                "SELECT * FROM users WHERE telegram_id = ?",
                (telegram_id,),
            )
            return cursor.fetchone()

    def set_user_sex(self, telegram_id: int, sex: str) -> None:
        """Set user's sex (male/female)."""
        with self.get_connection() as conn:
            conn.execute(
                "UPDATE users SET sex = ? WHERE telegram_id = ?",
                (sex, telegram_id),
            )
            conn.commit()

    def set_user_about(self, telegram_id: int, about: str) -> None:
        """Set user's about/introduction text."""
        with self.get_connection() as conn:
            conn.execute(
                "UPDATE users SET about = ? WHERE telegram_id = ?",
                (about, telegram_id),
            )
            conn.commit()

    def set_user_state(self, telegram_id: int, state: str) -> None:
        """Set user's current state in the flow."""
        with self.get_connection() as conn:
            conn.execute(
                "UPDATE users SET state = ? WHERE telegram_id = ?",
                (state, telegram_id),
            )
            conn.commit()

    def get_user_state(self, telegram_id: int) -> str:
        """Get user's current state."""
        with self.get_connection() as conn:
            cursor = conn.execute(
                "SELECT state FROM users WHERE telegram_id = ?",
                (telegram_id,),
            )
            row = cursor.fetchone()
            return row["state"] if row else "start"
