package repo_test

import (
	"context"
	"fmt"
	api "github.com/openshift-kni/oran-o2ims/internal/service/alarms/api/generated"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/db/models"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pashagolub/pgxmock/v4"

	alarmsrepo "github.com/openshift-kni/oran-o2ims/internal/service/alarms/internal/db/repo"
)

var _ = Describe("AlarmsRepository", func() {
	var (
		mock pgxmock.PgxPoolIface
		repo *alarmsrepo.AlarmsRepository
		ctx  context.Context
	)

	BeforeEach(func() {
		var err error
		mock, err = pgxmock.NewPool()
		Expect(err).NotTo(HaveOccurred())

		repo = &alarmsrepo.AlarmsRepository{
			Db: mock,
		}
		ctx = context.Background()
	})

	AfterEach(func() {
		mock.Close()
	})

	Describe("CreateServiceConfiguration", func() {
		When("no configuration exists", func() {
			It("creates a new configuration", func() {
				mock.ExpectQuery("SELECT (.+) FROM alarm_service_configuration").
					WillReturnRows(pgxmock.NewRows([]string{"id", "retention_period", "created_at", "updated_at"}))

				mock.ExpectQuery("INSERT INTO alarm_service_configuration").
					WithArgs(30).
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "retention_period", "created_at", "updated_at"}).
							AddRow(uuid.New(), 30, time.Now(), time.Now()),
					)

				config, err := repo.CreateServiceConfiguration(ctx, 30)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.RetentionPeriod).To(Equal(30))

				// Verify all expectations were met
				Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			})
		})
		When("one configuration exists", func() {
			It("returns the existing configuration as CreateServiceConfiguration is only called during init)", func() {
				existingID := uuid.New()
				mock.ExpectQuery("SELECT (.+) FROM alarm_service_configuration").
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "retention_period", "created_at", "updated_at"}).
							AddRow(existingID, 45, time.Now(), time.Now()),
					)

				config, err := repo.CreateServiceConfiguration(ctx, 30)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.RetentionPeriod).To(Equal(45)) // should return existing value, not input value
				Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			})
		})

		When("multiple configurations exist", func() {
			It("keeps first configuration and deletes others", func() {
				now := time.Now()
				id1, id2, id3 := uuid.New(), uuid.New(), uuid.New()

				mock.ExpectQuery("SELECT (.+) FROM alarm_service_configuration").
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "retention_period", "created_at", "updated_at"}).
							AddRow(id1, 45, now, now).
							AddRow(id2, 60, now, now).
							AddRow(id3, 90, now, now),
					)

				// Expect deletion of extra configurations (IDs 2 and 3)
				mock.ExpectExec("DELETE FROM alarm_definition WHERE").
					WithArgs(id2, id3).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))

				config, err := repo.CreateServiceConfiguration(ctx, 30)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.RetentionPeriod).To(Equal(45)) // should keep first config's value
				Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			})
		})

		When("delete operation fails", func() {
			It("returns an error", func() {
				now := time.Now()
				id1, id2 := uuid.New(), uuid.New()

				mock.ExpectQuery("SELECT (.+) FROM alarm_service_configuration").
					WillReturnRows(
						pgxmock.NewRows([]string{"id", "retention_period", "created_at", "updated_at"}).
							AddRow(id1, 45, now, now).
							AddRow(id2, 60, now, now),
					)

				// Simulate delete operation failure
				mock.ExpectExec("DELETE FROM alarm_definition WHERE").
					WithArgs(id2).
					WillReturnError(fmt.Errorf("database error"))

				config, err := repo.CreateServiceConfiguration(ctx, 30)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to delete additional service configurations"))
				Expect(config).To(BeNil())
				Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			})
		})

		When("initial query fails", func() {
			It("returns an error", func() {
				mock.ExpectQuery("SELECT (.+) FROM alarm_service_configuration").
					WillReturnError(fmt.Errorf("database error"))

				config, err := repo.CreateServiceConfiguration(ctx, 30)
				Expect(err).To(HaveOccurred())
				Expect(config).To(BeNil())
				Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
			})
		})
	})

	Describe("DeleteAlarmDefinitionsNotIn", func() {
		It("deletes alarm definitions not in provided IDs with correct object type ID", func() {
			ids := []any{uuid.New(), uuid.New()}
			objectTypeID := uuid.New()

			mock.ExpectExec("DELETE FROM alarm_definition WHERE").
				WithArgs(ids[0], ids[1], objectTypeID).
				WillReturnResult(pgxmock.NewResult("DELETE", 2))

			count, err := repo.DeleteAlarmDefinitionsNotIn(ctx, ids, objectTypeID)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))
		})

		It("returns error when deletion fails", func() {
			ids := []any{uuid.New()}
			objectTypeID := uuid.New()

			mock.ExpectExec("DELETE FROM alarm_definition WHERE").
				WithArgs(ids[0], objectTypeID).
				WillReturnError(fmt.Errorf("database error"))

			count, err := repo.DeleteAlarmDefinitionsNotIn(ctx, ids, objectTypeID)
			Expect(err).To(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})

		It("handles empty ID list", func() {
			var ids []any
			objectTypeID := uuid.New()

			mock.ExpectExec("DELETE FROM alarm_definition WHERE").
				WithArgs(objectTypeID).
				WillReturnResult(pgxmock.NewResult("DELETE", 0))

			count, err := repo.DeleteAlarmDefinitionsNotIn(ctx, ids, objectTypeID)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(0)))
		})
	})
	Describe("GetAlarmEventRecords", func() {
		When("records exist", func() {
			It("returns all alarm event records", func() {
				now := time.Now()
				id1, id2 := uuid.New(), uuid.New()

				mock.ExpectQuery("SELECT (.+) FROM alarm_event_record").
					WillReturnRows(
						pgxmock.NewRows([]string{
							"alarm_event_record_id", "alarm_raised_time", "alarm_acknowledged", "perceived_severity",
						}).
							AddRow(id1, now, false, api.CLEARED).
							AddRow(id2, now, true, api.INDETERMINATE),
					)

				records, err := repo.GetAlarmEventRecords(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(records).To(HaveLen(2))
				Expect(records[0].AlarmEventRecordID).To(Equal(id1))
				Expect(records[1].AlarmEventRecordID).To(Equal(id2))
			})
		})
	})

	Describe("PatchAlarmEventRecordACK", func() {
		When("the alarm exists", func() {
			It("updates acknowledgment status of an alarm", func() {
				id := uuid.New()
				now := time.Now()
				record := &models.AlarmEventRecord{
					AlarmEventRecordID:    id,
					AlarmAcknowledged:     true,
					AlarmAcknowledgedTime: &now,
					PerceivedSeverity:     api.WARNING,
				}

				mock.ExpectQuery("UPDATE alarm_event_record SET").
					WithArgs(true, &now, api.WARNING, id).
					WillReturnRows(
						pgxmock.NewRows([]string{
							"alarm_event_record_id", "alarm_acknowledged", "alarm_acknowledged_time", "alarm_cleared_time", "perceived_severity",
						}).AddRow(id, true, &now, &now, api.WARNING),
					)

				result, err := repo.PatchAlarmEventRecordACK(ctx, id, record)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.AlarmAcknowledged).To(BeTrue())
				Expect(result.AlarmAcknowledgedTime).To(Equal(&now))
			})
		})
	})

	Describe("GetAlarmSubscription", func() {
		When("subscription exists", func() {
			It("retrieves a specific alarm subscription", func() {
				id := uuid.New()
				conId := uuid.New()
				callback := "http://example.com/callback"

				mock.ExpectQuery("SELECT (.+) FROM alarm_subscription_info WHERE").
					WithArgs(id).
					WillReturnRows(
						pgxmock.NewRows([]string{
							"subscription_id", "consumer_subscription_id", "callback",
						}).AddRow(id, &conId, callback),
					)

				subscription, err := repo.GetAlarmSubscription(ctx, id)
				Expect(err).NotTo(HaveOccurred())
				Expect(subscription.SubscriptionID).To(Equal(id))
				Expect(subscription.ConsumerSubscriptionID).To(Equal(&conId))
				Expect(subscription.Callback).To(Equal(callback))
			})
		})

		When("subscription doesn't exist", func() {
			It("returns error", func() {
				id := uuid.New()
				mock.ExpectQuery("SELECT (.+) FROM alarm_subscription_info WHERE").
					WithArgs(id).
					WillReturnRows(pgxmock.NewRows([]string{}))

				subscription, err := repo.GetAlarmSubscription(ctx, id)
				Expect(err).To(HaveOccurred())
				Expect(subscription).To(BeNil())
			})
		})
	})

	Describe("UpsertAlarmEventRecord", func() {
		When("upserting a single record", func() {
			It("successfully upserts alarm event records", func() {
				id := uuid.New()
				records := []models.AlarmEventRecord{
					{
						AlarmRaisedTime:   time.Now(),
						PerceivedSeverity: api.WARNING,
						ObjectID:          &id,
						Fingerprint:       "fp1",
					},
				}

				mock.ExpectExec("INSERT INTO alarm_event_record").
					WithArgs(
						records[0].AlarmRaisedTime, records[0].AlarmClearedTime,
						records[0].AlarmAcknowledgedTime, records[0].AlarmAcknowledged,
						records[0].PerceivedSeverity, records[0].Extensions,
						records[0].ObjectID, records[0].ObjectTypeID,
						records[0].AlarmStatus, records[0].Fingerprint,
						records[0].AlarmDefinitionID, records[0].ProbableCauseID,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				err := repo.UpsertAlarmEventRecord(ctx, records)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("upserting multiple records", func() {
			It("handles multiple records in a single upsert", func() {
				id1, id2 := uuid.New(), uuid.New()
				now := time.Now()
				records := []models.AlarmEventRecord{
					{
						AlarmRaisedTime:   now,
						PerceivedSeverity: api.WARNING,
						ObjectID:          &id1,
						Fingerprint:       "fp1",
					},
					{
						AlarmRaisedTime:   now,
						PerceivedSeverity: api.CRITICAL,
						ObjectID:          &id2,
						Fingerprint:       "fp2",
					},
				}

				mock.ExpectExec("INSERT INTO alarm_event_record").
					WithArgs(
						records[0].AlarmRaisedTime, records[0].AlarmClearedTime,
						records[0].AlarmAcknowledgedTime, records[0].AlarmAcknowledged,
						records[0].PerceivedSeverity, records[0].Extensions,
						records[0].ObjectID, records[0].ObjectTypeID,
						records[0].AlarmStatus, records[0].Fingerprint,
						records[0].AlarmDefinitionID, records[0].ProbableCauseID,
						records[1].AlarmRaisedTime, records[1].AlarmClearedTime,
						records[1].AlarmAcknowledgedTime, records[1].AlarmAcknowledged,
						records[1].PerceivedSeverity, records[1].Extensions,
						records[1].ObjectID, records[1].ObjectTypeID,
						records[1].AlarmStatus, records[1].Fingerprint,
						records[1].AlarmDefinitionID, records[1].ProbableCauseID,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 2))

				err := repo.UpsertAlarmEventRecord(ctx, records)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("given an empty record list", func() {
			It("handles empty record list", func() {
				err := repo.UpsertAlarmEventRecord(ctx, []models.AlarmEventRecord{})
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("UpsertAlarmDefinitions", func() {
		When("upserting valid definitions", func() {
			It("successfully upserts alarm definitions", func() {
				defs := []models.AlarmDefinition{
					{
						AlarmName:         "test-alarm",
						AlarmLastChange:   "test",
						Severity:          string(api.AlarmSubscriptionInfoFilterNEW),
						AlarmDictionaryID: uuid.New(),
					},
				}

				mock.ExpectQuery("INSERT INTO alarm_definition").
					WithArgs(
						defs[0].AlarmName, defs[0].AlarmLastChange,
						defs[0].AlarmDescription, defs[0].ProposedRepairActions,
						defs[0].AlarmAdditionalFields, defs[0].AlarmDictionaryID,
						defs[0].Severity,
					).
					WillReturnRows(pgxmock.NewRows([]string{"alarm_definition_id"}).AddRow(uuid.New()))

				results, err := repo.UpsertAlarmDefinitions(ctx, defs)
				Expect(err).NotTo(HaveOccurred())
				Expect(results).To(HaveLen(1))
			})
		})

		When("upserting empty input", func() {
			It("handles empty input gracefully", func() {
				results, err := repo.UpsertAlarmDefinitions(ctx, []models.AlarmDefinition{})
				Expect(err).NotTo(HaveOccurred())
				Expect(results).To(HaveLen(0))
			})
		})
	})

	Describe("GetAlarmsForSubscription", func() {
		It("retrieves alarms based on subscription criteria", func() {
			f := api.AlarmSubscriptionInfoFilterNEW
			subscription := models.AlarmSubscription{
				SubscriptionID: uuid.New(),
				EventCursor:    5,
				Filter:         &f,
			}

			mock.ExpectQuery("SELECT (.+) FROM alarm_event_record WHERE").
				WithArgs(subscription.EventCursor, subscription.Filter).
				WillReturnRows(
					pgxmock.NewRows([]string{
						"alarm_event_record_id", "alarm_raised_time",
						"perceived_severity", "notification_event_type",
					}).
						AddRow(uuid.New(), time.Now(), api.WARNING, f),
				)

			results, err := repo.GetAlarmsForSubscription(ctx, subscription)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
		})

		It("filters alarms by notification event type", func() {
			f := api.AlarmSubscriptionInfoFilterACKNOWLEDGE
			subscription := models.AlarmSubscription{
				SubscriptionID: uuid.New(),
				EventCursor:    5,
				Filter:         &f,
			}

			mock.ExpectQuery("SELECT (.+) FROM alarm_event_record WHERE").
				WithArgs(subscription.EventCursor, subscription.Filter).
				WillReturnRows(
					pgxmock.NewRows([]string{
						"alarm_event_record_id", "alarm_raised_time",
						"perceived_severity", "notification_event_type",
						"alarm_sequence_number",
					}).AddRow(uuid.New(), time.Now(), api.WARNING, f, int64(6)))

			results, err := repo.GetAlarmsForSubscription(ctx, subscription)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].NotificationEventType).To(Equal(f))
		})

		It("handles subscription with no filter", func() {
			subscription := models.AlarmSubscription{
				SubscriptionID: uuid.New(),
				EventCursor:    5,
				Filter:         nil,
			}

			mock.ExpectQuery("SELECT (.+) FROM alarm_event_record WHERE").
				WithArgs(subscription.EventCursor).
				WillReturnRows(pgxmock.NewRows([]string{
					"alarm_event_record_id", "alarm_raised_time",
					"perceived_severity", "notification_event_type",
					"alarm_sequence_number",
				}))

			results, err := repo.GetAlarmsForSubscription(ctx, subscription)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(BeEmpty())
		})

		It("handles no alarms above event cursor", func() {
			subscription := models.AlarmSubscription{
				SubscriptionID: uuid.New(),
				EventCursor:    1000, // High cursor value
				Filter:         nil,
			}

			mock.ExpectQuery("SELECT (.+) FROM alarm_event_record WHERE").
				WithArgs(subscription.EventCursor).
				WillReturnRows(pgxmock.NewRows([]string{
					"alarm_event_record_id", "alarm_raised_time",
					"perceived_severity", "notification_event_type",
					"alarm_sequence_number",
				}))

			results, err := repo.GetAlarmsForSubscription(ctx, subscription)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(BeEmpty())
		})

	})

	Describe("DeleteAlarmDictionariesNotIn", func() {
		When("there are dictionaries to delete", func() {
			It("deletes alarm dictionaries not in provided IDs", func() {
				ids := []any{uuid.New(), uuid.New()}

				mock.ExpectExec("DELETE FROM alarm_dictionary WHERE").
					WithArgs(ids[0], ids[1]).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

				err := repo.DeleteAlarmDictionariesNotIn(ctx, ids)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("UpdateSubscriptionEventCursor", func() {
		When("update is successful", func() {
			It("updates subscription event cursor successfully", func() {
				subscription := models.AlarmSubscription{
					SubscriptionID: uuid.New(),
					EventCursor:    10,
				}

				mock.ExpectQuery("UPDATE alarm_subscription_info SET").
					WithArgs(subscription.EventCursor, subscription.SubscriptionID).
					WillReturnRows(pgxmock.NewRows([]string{"subscription_id"}).AddRow(subscription.SubscriptionID))

				err := repo.UpdateSubscriptionEventCursor(ctx, subscription)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("update fails", func() {
			It("returns error", func() {
				subscription := models.AlarmSubscription{
					SubscriptionID: uuid.New(),
					EventCursor:    10,
				}

				mock.ExpectQuery("UPDATE alarm_subscription_info SET").
					WithArgs(subscription.EventCursor, subscription.SubscriptionID).
					WillReturnError(fmt.Errorf("database error"))

				err := repo.UpdateSubscriptionEventCursor(ctx, subscription)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("UpsertAlarmDictionary", func() {
		When("upserting a new dictionary", func() {
			It("successfully upserts alarm dictionary", func() {
				record := models.AlarmDictionary{
					AlarmDictionaryVersion: "1.0",
					EntityType:             "TestEntity",
					Vendor:                 "TestVendor",
					ObjectTypeID:           uuid.New(),
				}

				mock.ExpectQuery("INSERT INTO alarm_dictionary").
					WithArgs(
						record.AlarmDictionaryVersion,
						record.EntityType,
						record.Vendor,
						record.ObjectTypeID,
					).
					WillReturnRows(pgxmock.NewRows([]string{"alarm_dictionary_id"}).AddRow(uuid.New()))

				results, err := repo.UpsertAlarmDictionary(ctx, record)
				Expect(err).NotTo(HaveOccurred())
				Expect(results).To(HaveLen(1))
			})
		})
	})

	Describe("ResolveNotificationIfNotInCurrent", func() {
		When("resolving notifications", func() {
			It("resolves notifications not in current payload", func() {
				clearTime := time.Now()
				alarmsrepo.TimeNow = func() time.Time {
					return clearTime
				}
				fp := "9a9e2d82a78cf2b9"
				t := time.Now()
				am := &api.AlertmanagerNotification{
					Alerts: []api.Alert{
						{
							Fingerprint: &fp,
							StartsAt:    &t,
						},
					},
				}
				mock.ExpectQuery("UPDATE alarm_event_record SET").
					WithArgs(api.Resolved, clearTime, api.CLEARED, &fp, &t).
					WillReturnRows(pgxmock.NewRows([]string{"alarm_event_record_id"}).AddRow(uuid.New()))

				err := repo.ResolveNotificationIfNotInCurrent(ctx, am)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("GetAlarmDefinitions", func() {
		When("valid alertmanager notification is provided", func() {
			It("retrieves alarm definitions for alertmanager notification", func() {
				fp := "9a9e2d82a78cf2b9"
				clusterID := uuid.New()
				alertname := "CollectorNodeDown"
				severity := "info"
				am := &api.AlertmanagerNotification{
					Alerts: []api.Alert{
						{
							Annotations: nil,
							Fingerprint: &fp,
							Labels: &map[string]string{
								"alertname":       alertname,
								"severity":        severity,
								"managed_cluster": clusterID.String(),
							},
							StartsAt: nil,
							Status:   nil,
						},
					},
				}

				objectTypeID := uuid.New()
				clusterMap := map[uuid.UUID]uuid.UUID{clusterID: objectTypeID}

				mock.ExpectQuery("SELECT (.+) FROM alarm_definition WHERE").
					WithArgs(alertname, objectTypeID, severity).
					WillReturnRows(pgxmock.NewRows([]string{
						"alarm_name", "alarm_definition_id", "probable_cause_id",
						"object_type_id", "severity",
					}).AddRow(
						alertname, uuid.New(), uuid.New(),
						objectTypeID, severity,
					))

				results, err := repo.GetAlarmDefinitions(ctx, am, clusterMap)
				Expect(err).NotTo(HaveOccurred())
				Expect(results).To(HaveLen(1))
			})
		})
	})

	Describe("DeleteAlarmSubscription", func() {
		When("subscription exists", func() {
			It("deletes subscription successfully", func() {
				id := uuid.New()

				mock.ExpectExec("DELETE FROM alarm_subscription_info WHERE").
					WithArgs(id).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

				count, err := repo.DeleteAlarmSubscription(ctx, id)
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(int64(1)))
			})
		})
	})

	Describe("CreateAlarmSubscription", func() {
		It("creates new subscription successfully", func() {
			f := api.AlarmSubscriptionInfoFilterACKNOWLEDGE
			id := uuid.New()
			subscription := models.AlarmSubscription{
				ConsumerSubscriptionID: &id,
				Callback:               "http://test.com/callback",
				Filter:                 &f,
				EventCursor:            int64(10),
			}

			mock.ExpectQuery("INSERT INTO alarm_subscription_info").
				WithArgs(
					subscription.ConsumerSubscriptionID, subscription.Filter, subscription.Callback,
					subscription.EventCursor,
				).
				WillReturnRows(pgxmock.NewRows([]string{"updated_at", "subscription_id", "consumer_subscription_id", "filter", "callback", "event_cursor", "created_at"}).
					AddRow(time.Now(), uuid.New(), &id, &f, "http://test.com/callback", int64(10), time.Now()))

			result, err := repo.CreateAlarmSubscription(ctx, subscription)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.ConsumerSubscriptionID).To(Equal(subscription.ConsumerSubscriptionID))
		})
	})

	Describe("GetMaxAlarmSeq", func() {
		When("alarms exist", func() {
			It("returns maximum alarm sequence number", func() {
				mock.ExpectQuery(`SELECT (.+MAX.+) FROM "alarm_event_record"`).
					WillReturnRows(pgxmock.NewRows([]string{"coalesce"}).AddRow(int64(42)))

				maxSeq, err := repo.GetMaxAlarmSeq(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(maxSeq).To(Equal(int64(42)))
			})
		})

		When("no alarms exist", func() {
			It("returns 0 when no alarms exist", func() {
				mock.ExpectQuery("SELECT COALESCE").
					WillReturnRows(pgxmock.NewRows([]string{"coalesce"}).AddRow(int64(0)))

				maxSeq, err := repo.GetMaxAlarmSeq(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(maxSeq).To(Equal(int64(0)))
			})
		})
	})
})
