import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { apiService } from '../services/api';
import './EmailSettings.css';
import axios from 'axios';

interface EmailSchedule {
  id: number;
  account_id: number;
  email_type: string;
  frequency: string;
  day_of_week?: number;
  day_of_month?: number;
  time_of_day: string;
  is_active: boolean;
  last_sent_at?: string;
  created_at: string;
  updated_at: string;
}

interface EmailScheduleRequest {
  email_type: string;
  frequency: string;
  day_of_week?: number;
  day_of_month?: number;
  time_of_day: string;
  is_active: boolean;
}

const EmailSettings: React.FC = () => {
  const { user } = useAuth();
  const [schedules, setSchedules] = useState<EmailSchedule[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingSchedule, setEditingSchedule] = useState<EmailSchedule | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  const emailTypes = [
    { value: 'weekly_stock_report', label: 'Weekly Stock Report' },
    { value: 'weekly_supply_chain_report', label: 'Weekly Supply Chain Report' },
    { value: 'low_stock_alert', label: 'Low Stock Alert' }
  ];

  const daysOfWeek = [
    { value: 0, label: 'Sunday' },
    { value: 1, label: 'Monday' },
    { value: 2, label: 'Tuesday' },
    { value: 3, label: 'Wednesday' },
    { value: 4, label: 'Thursday' },
    { value: 5, label: 'Friday' },
    { value: 6, label: 'Saturday' }
  ];

  const frequencies = [
    { value: 'weekly', label: 'Weekly' },
    { value: 'daily', label: 'Daily' },
    { value: 'monthly', label: 'Monthly' }
  ];

  useEffect(() => {
    if (user?.account_id) {
      fetchSchedules();
    }
  }, [user?.account_id]);

  const fetchSchedules = async () => {
    try {
      setLoading(true);
      const response = await axios.get(`/accounts/${user?.account_id}/email-schedules`);
      setSchedules(response.data.schedules);
      setError(null);
    } catch (err) {
      setError('Failed to load email schedules');
      console.error('Error fetching schedules:', err);
    } finally {
      setLoading(false);
    }
  };

  const toggleSchedule = async (emailType: string) => {
    try {
      const response = await axios.patch(`/accounts/${user?.account_id}/email-schedules/${emailType}/toggle`);
      setSchedules(prev => prev.map(schedule => 
        schedule.email_type === emailType 
          ? { ...schedule, is_active: response.data.is_active }
          : schedule
      ));
    } catch (err) {
      setError('Failed to toggle email schedule');
      console.error('Error toggling schedule:', err);
    }
  };

  const deleteSchedule = async (emailType: string) => {
    if (!window.confirm('Are you sure you want to delete this email schedule?')) {
      return;
    }

    try {
      await axios.delete(`/accounts/${user?.account_id}/email-schedules/${emailType}`);
      setSchedules(prev => prev.filter(schedule => schedule.email_type !== emailType));
    } catch (err) {
      setError('Failed to delete email schedule');
      console.error('Error deleting schedule:', err);
    }
  };

  const createSchedule = async (scheduleData: EmailScheduleRequest) => {
    try {
      const response = await axios.post(`/accounts/${user?.account_id}/email-schedules`, scheduleData);
      setSchedules(prev => [...prev, response.data]);
      setShowCreateForm(false);
    } catch (err) {
      setError('Failed to create email schedule');
      console.error('Error creating schedule:', err);
    }
  };

  const updateSchedule = async (emailType: string, scheduleData: EmailScheduleRequest) => {
    try {
      const response = await axios.put(`/accounts/${user?.account_id}/email-schedules/${emailType}`, scheduleData);
      setSchedules(prev => prev.map(schedule => 
        schedule.email_type === emailType ? response.data : schedule
      ));
      setEditingSchedule(null);
    } catch (err) {
      setError('Failed to update email schedule');
      console.error('Error updating schedule:', err);
    }
  };

  const getEmailTypeLabel = (emailType: string) => {
    return emailTypes.find(type => type.value === emailType)?.label || emailType;
  };

  const getDayLabel = (dayOfWeek?: number) => {
    if (dayOfWeek === undefined) return 'Not set';
    return daysOfWeek.find(day => day.value === dayOfWeek)?.label || 'Unknown';
  };

  const getFrequencyLabel = (frequency: string) => {
    return frequencies.find(freq => freq.value === frequency)?.label || frequency;
  };

  if (loading) {
    return <div className="email-settings-loading">Loading email settings...</div>;
  }

  return (
    <div className="email-settings">
      <div className="email-settings-header">
        <h2>Email Settings</h2>
        <button 
          className="btn btn-primary"
          onClick={() => setShowCreateForm(true)}
        >
          Add Email Schedule
        </button>
      </div>

      {error && (
        <div className="error-message">
          {error}
          <button onClick={() => setError(null)}>×</button>
        </div>
      )}

      <div className="email-schedules-list">
        {schedules.length === 0 ? (
          <div className="no-schedules">
            <p>No email schedules configured. Create one to get started.</p>
          </div>
        ) : (
          schedules.map(schedule => (
            <div key={schedule.id} className="schedule-card">
              <div className="schedule-header">
                <h3>{getEmailTypeLabel(schedule.email_type)}</h3>
                <div className="schedule-status">
                  <span className={`status-badge ${schedule.is_active ? 'active' : 'inactive'}`}>
                    {schedule.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
              </div>

              <div className="schedule-details">
                <p><strong>Frequency:</strong> {getFrequencyLabel(schedule.frequency)}</p>
                {schedule.day_of_week !== undefined && (
                  <p><strong>Day:</strong> {getDayLabel(schedule.day_of_week)}</p>
                )}
                <p><strong>Time:</strong> {schedule.time_of_day}</p>
                {schedule.last_sent_at && (
                  <p><strong>Last Sent:</strong> {new Date(schedule.last_sent_at).toLocaleString()}</p>
                )}
              </div>

              <div className="schedule-actions">
                <button 
                  className="btn btn-secondary"
                  onClick={() => toggleSchedule(schedule.email_type)}
                >
                  {schedule.is_active ? 'Disable' : 'Enable'}
                </button>
                <button 
                  className="btn btn-secondary"
                  onClick={() => setEditingSchedule(schedule)}
                >
                  Edit
                </button>
                <button 
                  className="btn btn-danger"
                  onClick={() => deleteSchedule(schedule.email_type)}
                >
                  Delete
                </button>
              </div>
            </div>
          ))
        )}
      </div>

      {showCreateForm && (
        <EmailScheduleForm
          onSubmit={createSchedule}
          onCancel={() => setShowCreateForm(false)}
          emailTypes={emailTypes}
          daysOfWeek={daysOfWeek}
          frequencies={frequencies}
        />
      )}

      {editingSchedule && (
        <EmailScheduleForm
          schedule={editingSchedule}
          onSubmit={(data) => updateSchedule(editingSchedule.email_type, data)}
          onCancel={() => setEditingSchedule(null)}
          emailTypes={emailTypes}
          daysOfWeek={daysOfWeek}
          frequencies={frequencies}
        />
      )}
    </div>
  );
};

interface EmailScheduleFormProps {
  schedule?: EmailSchedule;
  onSubmit: (data: EmailScheduleRequest) => void;
  onCancel: () => void;
  emailTypes: Array<{ value: string; label: string }>;
  daysOfWeek: Array<{ value: number; label: string }>;
  frequencies: Array<{ value: string; label: string }>;
}

const EmailScheduleForm: React.FC<EmailScheduleFormProps> = ({
  schedule,
  onSubmit,
  onCancel,
  emailTypes,
  daysOfWeek,
  frequencies
}) => {
  const [formData, setFormData] = useState<EmailScheduleRequest>({
    email_type: schedule?.email_type || 'weekly_stock_report',
    frequency: schedule?.frequency || 'weekly',
    day_of_week: schedule?.day_of_week,
    day_of_month: schedule?.day_of_month,
    time_of_day: schedule?.time_of_day || '09:00',
    is_active: schedule?.is_active ?? true
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <div className="email-schedule-form-overlay">
      <div className="email-schedule-form">
        <h3>{schedule ? 'Edit Email Schedule' : 'Create Email Schedule'}</h3>
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="email_type">Email Type</label>
            <select
              id="email_type"
              value={formData.email_type}
              onChange={(e) => setFormData({ ...formData, email_type: e.target.value })}
              disabled={!!schedule} // Can't change email type when editing
            >
              {emailTypes.map(type => (
                <option key={type.value} value={type.value}>
                  {type.label}
                </option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="frequency">Frequency</label>
            <select
              id="frequency"
              value={formData.frequency}
              onChange={(e) => setFormData({ ...formData, frequency: e.target.value })}
            >
              {frequencies.map(freq => (
                <option key={freq.value} value={freq.value}>
                  {freq.label}
                </option>
              ))}
            </select>
          </div>

          {formData.frequency === 'weekly' && (
            <div className="form-group">
              <label htmlFor="day_of_week">Day of Week</label>
              <select
                id="day_of_week"
                value={formData.day_of_week || ''}
                onChange={(e) => setFormData({ 
                  ...formData, 
                  day_of_week: e.target.value ? parseInt(e.target.value) : undefined 
                })}
              >
                <option value="">Select day</option>
                {daysOfWeek.map(day => (
                  <option key={day.value} value={day.value}>
                    {day.label}
                  </option>
                ))}
              </select>
            </div>
          )}

          <div className="form-group">
            <label htmlFor="time_of_day">Time of Day</label>
            <input
              type="time"
              id="time_of_day"
              value={formData.time_of_day}
              onChange={(e) => setFormData({ ...formData, time_of_day: e.target.value })}
            />
          </div>

          <div className="form-group">
            <label>
              <input
                type="checkbox"
                checked={formData.is_active}
                onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
              />
              Active
            </label>
          </div>

          <div className="form-actions">
            <button type="submit" className="btn btn-primary">
              {schedule ? 'Update' : 'Create'}
            </button>
            <button type="button" className="btn btn-secondary" onClick={onCancel}>
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EmailSettings; 