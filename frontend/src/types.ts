export interface Market {
  id: string;
  title: string;
  status: 'DRAFT' | 'OPEN' | 'LOCKED' | 'EVENT_ACTIVE' | 'RESOLUTION_PENDING' | 'RESOLVED' | 'VOID' | 'SETTLED' | 'REFUNDED';
  event_window_seconds: number;
  resolution_method: string;
  rules: string;
}

export interface MarketTemplate {
  id: string;
  title_pattern: string;
  resolution_method: string;
  default_rules: string;
}
