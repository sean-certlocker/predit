export interface Market {
  id: string;
  title: string;
  status: 'DRAFT' | 'OPEN' | 'LOCKED' | 'EVENT_ACTIVE' | 'RESOLUTION_PENDING' | 'RESOLVED' | 'VOID' | 'SETTLED' | 'REFUNDED';
  event_window_seconds: number;
  resolution_method: string;
  rules: string;
  health?: string;
  safety_flags?: string[];
}

export interface User {
  id: string;
  username: string;
  role: string;
  trust_score: number;
  suspicious: boolean;
}

export interface MarketTemplate {
  id: string;
  title_pattern: string;
  resolution_method: string;
  default_rules: string;
}
