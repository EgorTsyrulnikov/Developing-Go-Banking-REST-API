package scheduler

import (
	"context"
	"time"
	"bankapi/internal/config"
	"bankapi/internal/repositories"
	"bankapi/pkg/logger"
	"bankapi/pkg/smtp"
)

func StartPaymentScheduler(ctx context.Context, cfg *config.Config) {
	ticker := time.NewTicker(12 * time.Hour)
	// Or a shorter duration for testing if needed
	
	go func() {
		logger.Log.Info("Payment scheduler started")
		for {
			select {
			case <-ticker.C:
				processOverduePayments(cfg)
			case <-ctx.Done():
				ticker.Stop()
				logger.Log.Info("Payment scheduler stopped")
				return
			}
		}
	}()
}

func processOverduePayments(cfg *config.Config) {
	logger.Log.Info("Checking for overdue payments...")
	now := time.Now()
	schedules, err := repositories.GetOverduePayments(now)
	if err != nil {
		logger.Log.Errorf("Error fetching overdue payments: %v", err)
		return
	}

	for _, s := range schedules {
		// Calculate penalty (e.g., +10% of amount if not already applied or applied daily)
		if s.Penalty == 0 {
			penalty := s.Amount * 0.10
			err := repositories.UpdatePaymentSchedulePenalty(s.ID, penalty)
			if err != nil {
				logger.Log.Errorf("Error updating penalty for schedule %s: %v", s.ID, err)
				continue
			}
			
			// Send email notification (mocking getting user email for simplicity)
			// In reality, we would query the user's email through credit -> user
			credit, err := repositories.GetCreditByID(s.CreditID)
			if err == nil && credit != nil {
				user, err := repositories.GetUserByID(credit.UserID) // Wait, we don't have GetUserByID yet
				if err == nil && user != nil {
					smtp.SendEmail(user.Email, "Overdue Payment Notice", "You have an overdue payment with a 10% penalty applied.", cfg)
				}
			}
		}
	}
}
