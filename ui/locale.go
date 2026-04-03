package ui

// locale.go — централизованный слой русской локализации.
//
// Все пользовательские строки UI хранятся здесь.
// Не меняй: package names, module path, binary name, import paths,
// CLI flags, env var names, config keys, внутренние slug/ID.

// ── Общие элементы управления ─────────────────────────────────────────────

const (
	strConfirm  = "подтвердить"
	strCancel   = "отмена"
	strClose    = "закрыть"
	strDelete   = "удалить"
	strNavigate = "навигация"
	strJoin     = "войти"
	strPost     = "опубликовать"
	strSubmitBtn = "отправить"
	strNextBtn  = "далее"
	strPrevNext = "предыдущее/следующее"
	strJump     = "перейти"
	strBrowse   = "просмотр"
	strSend     = "отправить"
	strSelect   = "выбрать"
	strVote     = "голосовать"
	strSignUp   = "записаться"
	strContinue = "продолжить"
	strNext     = "далее"
	strLatest   = " (последняя)"
)

// ── Справка (Help Modal) ───────────────────────────────────────────────────

const (
	strHelpTitle       = " Помощь "
	strHelpCatChat     = "КЛАВИШИ ЧАТА"
	strHelpCatCommands = "КОМАНДЫ"
	strHelpCatGallery  = "КЛАВИШИ ГАЛЕРЕИ"
	strHelpCatInfo     = "ИНФОРМАЦИЯ"
	strHelpInfoLine1   = "  Все данные сбрасываются каждое воскресенье в 23:59 UTC."
	strHelpInfoLine2   = "  Ничто не вечно. Пиши пока можешь."

	// Клавиши чата
	strKeyHelpThis   = "эта справка"
	strKeyChangeNick = "сменить ник"
	strKeySwitchRooms = "сменить комнату"
	strKeyViewMentions = "упоминания"
	strKeyPostNote   = "написать заметку"
	strKeyTankard    = "кружка"
	strKeyLeaderboard = "рейтинг"
	strKeyToggleDMs  = "переключить ЛС"
	strKeyCloseModal = "закрыть / снять фокус"
	strKeyScrollChat = "прокрутить чат"

	// Команды
	strCmdPoll        = "создать опрос"
	strCmdVote        = "голосовать в опросе"
	strCmdEndpoll     = "закрыть свой опрос"
	strCmdGif         = "найти и отправить GIF"
	strCmdSubmit      = "ввести флаг варгейма"
	strCmdLeaderboard = "рейтинг хакеров"
	strCmdDm          = "открыть личку"
	strCmdDmMsg       = "отправить ЛС напрямую"

	// Клавиши галереи
	strGalleryKeyPost   = "написать заметку"
	strGalleryKeyDelete = "удалить заметку"
	strGalleryKeyTab    = "выбор"
	strGalleryKeyDrag   = "переместить заметку"
)

// ── Смена ника (Nick Modal) ────────────────────────────────────────────────

const (
	strNickTitle       = " Сменить ник "
	strNickHint        = "  Введите новый ник (2-20 символов)"
	strNickPlaceholder = "новый ник..."
	strNickErrLen      = "ник должен быть 2-20 символов"
)

// ── Выбор комнаты (Join Room Modal) ───────────────────────────────────────

const (
	strJoinRoomTitle  = " Выбрать комнату "
	strJoinRoomHere   = "(здесь)"
	strLabelRooms     = "  КОМНАТЫ"
	strLabelWargames  = "  ВАРГЕЙМЫ"
)

// ── Написать заметку (Post Note Modal) ────────────────────────────────────

const (
	strPostNoteTitle       = " Написать заметку "
	strPostNoteHint        = "  Максимум 280 символов. Текст переносится."
	strPostNotePlaceholder = "напишите что-нибудь на доске..."
)

// ── Просмотр заметки (Expand Note Modal) ──────────────────────────────────

const (
	strNoteTitle = " Заметка "
)

// ── Подтверждение (Admin Confirm Modal) ───────────────────────────────────

const (
	strAdminTitle     = " АДМИН "
	strAdminCannotUndo = "Это действие нельзя отменить."
)

// ── Упоминание (Mention Modal) ─────────────────────────────────────────────

// strMentionHeaderFmt форматируется с автором и комнатой.
const strMentionHeaderFmt = " упоминание от %s в #%s "

// ── Варгейм: правила (Wargame Rules Modal) ────────────────────────────────

const (
	strWargameCatWhat    = "ЧТО ЭТО"
	strWargameWhat1      = "  Практикуй хакерские задания из"
	strWargameWhat2      = "  OverTheWire и других варгеймов."
	strWargameWhat3      = "  Решай уровни, вводи флаги, зарабатывай"
	strWargameWhat4      = "  очки и поднимайся в рейтинге."
	strWargameCatHow     = "КАК ИГРАТЬ"
	strWargameHow1Pre    = "  1. Зарегистрируйся с "
	strWargameHow1Suf    = " ниже"
	strWargameHow2Pre    = "  2. Перейди на "
	strWargameHow3       = "  3. Реши уровни, найди флаг"
	strWargameHow4Pre    = "  4. Введи "
	strWargameHow4Suf    = " чтобы отправить"
	strWargameHow5       = "  5. Зарабатывай очки и поднимайся"
	strWargameCatPoints  = "ОЧКИ"
	strWargamePoints1    = "  Уровень N = N очков"
	strWargamePoints2    = "  Ур.1=1  Ур.5=5  Ур.10=10"
	strWargamePoints3    = "  Очки сохраняются навсегда."
	strWargameSignedUp   = "  В ИГРЕ"
	strWargameSignedUpSuf = " — ты в рейтинге"
	strWargameNotSigned  = "  НЕ ЗАРЕГИСТРИРОВАН"
)

// ── Рейтинг (Leaderboard Modal) ───────────────────────────────────────────

const (
	strLeaderboardTitle  = "РЕЙТИНГ"
	strLeaderboardEmpty  = "  Пока никого в рейтинге."
	strProgressTitle     = "ВАШ ПРОГРЕСС"
	strProgressNoFlags   = "нет флагов"
	// strProgressTotalFmt форматируется с level, points, rank.
	strProgressTotalFmt = "  Итого: Ур.%d  %d очков  Место: %s"
	strProgressUnranked  = "нет места"
)

// ── Ввод флага (Submit Flag Modal) ────────────────────────────────────────

const (
	strSubmitTitle       = "ВВОД ФЛАГА"
	strSubmitLevelFmt    = "  Уровень %d"
	strSubmitFlagLabel   = "Флаг:"
	strSubmitFlagEmpty   = "флаг не может быть пустым"
	strSubmitPlaceholder = "вставьте флаг..."
)

// ── Создать опрос (Poll Creation Modal) ───────────────────────────────────

const (
	strPollTitle            = " Создать опрос "
	strPollLabelTitle       = "Заголовок"
	strPollOptionFmt        = "Вариант %d"
	strPollAddOptionHint    = "  ctrl+n чтобы добавить вариант"
	strPollErrTitleRequired = "необходим заголовок"
	strPollErrMinOptions    = "нужно минимум 2 варианта"
	strPollQuestionPH       = "вопрос опроса..."
	strPollOptionPHFmt      = "вариант %d..."
)

// ── Голосование (Poll Vote Overlay) ───────────────────────────────────────

const (
	// strPollVoteHeaderFmt: current, total
	strPollVoteHeaderFmt       = " ОПРОСЫ [%d/%d] "
	strPollVoteHeaderClosedFmt = " ОПРОСЫ [%d/%d] (закрыт) "
	// strPollVoteByFmt: creatorNick, totalVotes
	strPollVoteByFmt  = "от %s · %d голосов"
	strPollCardOpen   = "ОПРОС"
	strPollCardClosed = "ОПРОС (закрыт)"
	// strPollVotesFmt: totalVotes
	strPollVotesFmt = "%d голосов"
	// strPollCastFmt: totalVotes
	strPollCastFmt = "/vote чтобы голосовать · %d голосов"
)

// ── GIF Поиск (Gif Modal) ─────────────────────────────────────────────────

const (
	// strGifTitleFmt: query
	strGifTitleFmt = " Поиск KLIPY: %s "
	strGifLoading  = "  загрузка..."
	strGifNoResult = "нет результатов"
)

// ── Журнал изменений (Changelog Modal) ────────────────────────────────────

const (
	strChangelogTitle = " История изменений "
)

// ── Заставка (Splash Screen) ─────────────────────────────────────────────

const (
	strSplashYouAre    = "вы — "
	strSplashLine1     = "терминальная таверна через SSH."
	strSplashLine2     = "чат, игры, общение."
	strSplashLine3     = "без аккаунтов. без логов. без правил."
	strSplashLine4     = "говори что думаешь."
	strSplashLine5     = "всё сбрасывается каждое воскресенье."
	strSplashLine6     = "ничто не вечно."
	strSplashEnterDesc = " войти в таверну"
	strSplashQuitDesc  = " выход"
	strSplashChangelog = " история изменений"
)

// ── Нижняя панель (Bottom Bar) ────────────────────────────────────────────

const (
	strBBarHelp        = "помощь"
	strBBarNick        = "ник"
	strBBarRooms       = "комнаты"
	strBBarMentions    = "упоминания"
	strBBarMentionsFmt = "упоминания(%d)"
	strBBarDMs         = "ЛС"
	strBBarTankard     = "кружка"
	strBBarLeaderboard = "рейтинг"
	strBBarScroll      = "прокрутка"
	strBBarTavern      = "таверна"
	strBBarBack        = "назад"
	strBBarNavigate    = "навигация"
	strBBarOpen        = "открыть"
	strBBarDrink       = "выпить"
	strBBarPost        = "написать"
	strBBarExpand      = "развернуть"
	strBBarDelete      = "удалить"
	strBBarSelect      = "выбор"
)

// ── Боковая панель (Sidebar) ──────────────────────────────────────────────

const (
	strSidebarRooms    = "КОМНАТЫ"
	strSidebarWargames = "ВАРГЕЙМЫ"
	strSidebarDMs      = "ЛС"
	strSidebarDMHint   = " TAB — открыть"
	strSidebarOtherSSH = "ДРУГИЕ SSH"
	strSidebarOnline   = "В СЕТИ"
	strSidebarEmpty    = "(нет)"
	strSidebarHackers  = "ХАКЕРЫ"
)

// ── Верхняя панель (Top Bar) ──────────────────────────────────────────────

const (
	strTopOnlineFmt  = "%d онлайн"
	strTopWeekFmt    = "%d на неделе"
	strTopAllTimeFmt = "%d всего"
)

// ── Чат (Chat View) ───────────────────────────────────────────────────────

const (
	strChatPlaceholder = "Введите сообщение..."
	// Индикаторы набора
	strTypingOne  = "%s печатает%s"
	strTypingTwo  = "%s и %s печатают%s"
	strTypingMany = "%d человек(а) печатают%s"
	// Временны́е метки
	strTsJustNow = "только что"
	strTsSecsAgo = "%dс назад"
	strTs1MinAgo = "1 мин назад"
	strTsMinsAgo = "%d мин назад"
)

// ── Личные сообщения: входящие (DM Inbox) ────────────────────────────────

const (
	strDMTitle        = "ЛИЧНЫЕ СООБЩЕНИЯ"
	strDMEmpty        = "Нет переписок."
	strDMEmptyHintPre = "Используй "
	strDMEmptyHintPos = " в таверне чтобы начать."
	strDMFooter       = "  ↑↓ навигация · ENTER открыть · TAB назад в таверну"
	strDMNowFmt       = "сейчас"
	strDMNewFmt       = " %d новых"
)

// ── Системные сообщения (app.go) ──────────────────────────────────────────

const (
	strSysWelcome      = "Добро пожаловать в таверну. /help — список команд."
	strSysWrongFlag    = "Неверный флаг. Попробуй ещё раз."
	// strSysHackedFmt: nick, room, level, totalLevel, totalPoints
	strSysHackedFmt    = ">> %s взломал %s уровень %d  [Ур.%d | %d очков]"
	// strSysJoinWargamesFmt: nick
	strSysJoinWargamesFmt = "%s вступил в варгейм"
	// strSysJoinedRoomFmt: room name
	strSysJoinedRoomFmt   = "Вы вошли в #%s"
	// strSysJoinedFromFmt: nick, old room
	strSysJoinedFromFmt   = "%s перешёл из #%s"
)

// ── Системные сообщения (internal/server/server.go) ──────────────────────
// Примечание: эти строки не используются здесь напрямую —
// они вынесены для документации. Переводы применены инлайн в server.go.

const (
	// strSrvJoinedFmt: nickname
	strSrvJoinedFmt = "%s зашёл в таверну"
	// strSrvLeftFmt: nickname
	strSrvLeftFmt = "%s покинул таверну"
)
