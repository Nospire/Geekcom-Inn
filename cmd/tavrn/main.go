package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tavrn/internal/admin"
	"tavrn/internal/hub"
	"tavrn/internal/server"
	"tavrn/internal/session"
	"tavrn/internal/store"
)

func main() {
	st, err := store.New("tavrn.db")
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	defer st.Close()

	h := hub.New()
	go h.Run()

	adminFP := os.Getenv("TAVRN_ADMIN_FP")
	adm := admin.New(st, adminFP)

	if _, err := os.Stat(".ssh"); os.IsNotExist(err) {
		os.MkdirAll(".ssh", 0700)
	}

	srv, err := server.New(server.Config{
		Host:             "0.0.0.0",
		Port:             2222,
		HostKeyPath:      ".ssh/id_ed25519",
		AdminFingerprint: adminFP,
		Store:            st,
		Hub:              h,
		Admin:            adm,
	})
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	// HTTP kill switch
	adminToken := os.Getenv("TAVRN_ADMIN_TOKEN")
	if adminToken != "" {
		srv.StartHTTPAdmin("127.0.0.1:8080", adminToken)
	}

	// Weekly purge scheduler
	go startPurgeScheduler(st, h)

	// Start SSH server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	log.Println("tavrn.sh is open. ssh localhost -p 2222")

	<-done
	log.Println("tavern closing...")
	h.BroadcastAll(session.Msg{
		Type: session.MsgSystem,
		Text: "the tavern is closing...",
	})
	srv.Shutdown(5 * time.Second)
	log.Println("goodbye.")
}

func startPurgeScheduler(st *store.Store, h *hub.Hub) {
	for {
		now := time.Now().UTC()
		daysUntilSunday := (7 - int(now.Weekday())) % 7
		if daysUntilSunday == 0 && (now.Hour() > 23 || (now.Hour() == 23 && now.Minute() >= 59)) {
			daysUntilSunday = 7
		}
		next := time.Date(now.Year(), now.Month(), now.Day()+daysUntilSunday, 23, 59, 0, 0, time.UTC)
		timer := time.NewTimer(time.Until(next))
		<-timer.C

		log.Println("Weekly purge starting...")
		h.BroadcastAll(session.Msg{
			Type: session.MsgSystem,
			Text: "The tavern has been swept clean.",
		})
		st.PurgeAll()
		log.Println("Weekly purge complete")
	}
}
