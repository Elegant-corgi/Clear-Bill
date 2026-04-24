export type ActivityStatus = "success" | "warning" | "processing";
export type BillStatus = "ready" | "reviewing" | "overdue" | "paid";
export type ReconciliationStatus = "running" | "done" | "attention";

export interface ActivityItem {
  title: string;
  description: string;
  time: string;
  status: ActivityStatus;
}

export interface AttentionItem {
  title: string;
  owner: string;
  dueDate: string;
  progress: number;
}

export interface DashboardSummary {
  openBills: number;
  pendingInvoices: number;
  overdueAmount: number;
  autoMatchedRate: number;
  recentActivity: ActivityItem[];
  attentionList: AttentionItem[];
}

export interface BillRecord {
  id: string;
  customerName: string;
  billingMonth: string;
  amount: number;
  status: BillStatus;
  channel: string;
  updatedAt: string;
}

export interface CustomerSnapshot {
  id: string;
  name: string;
  creditLevel: string;
  activeContracts: number;
  outstandingAmount: number;
  billingHealth: number;
  primaryContact: string;
}

export interface ReconciliationTask {
  id: string;
  bank: string;
  period: string;
  matched: number;
  total: number;
  status: ReconciliationStatus;
  owner: string;
}
