import type {
  BillRecord,
  CustomerSnapshot,
  DashboardSummary,
  ReconciliationTask,
} from "./types";

export const dashboardSummaryMock: DashboardSummary = {
  openBills: 128,
  pendingInvoices: 37,
  overdueAmount: 482130.5,
  autoMatchedRate: 92.4,
  recentActivity: [
    {
      title: "Northwind April bill imported",
      description: "4,216 line items were normalized and queued for review.",
      time: "09:20",
      status: "success",
    },
    {
      title: "2 reconciliation groups need confirmation",
      description: "Bank feed and invoice ledger mismatch exceeded tolerance.",
      time: "08:40",
      status: "warning",
    },
    {
      title: "Auto-posting window started",
      description: "Scheduled posting is processing approved invoices.",
      time: "08:00",
      status: "processing",
    },
  ],
  attentionList: [
    {
      title: "Government account overdue review",
      owner: "Avery",
      dueDate: "Today 15:00",
      progress: 68,
    },
    {
      title: "Q2 pre-close billing lock",
      owner: "Nina",
      dueDate: "Tomorrow 10:30",
      progress: 44,
    },
    {
      title: "Tax inclusive invoice spot check",
      owner: "Leo",
      dueDate: "Apr 28",
      progress: 82,
    },
  ],
};

export const billRecordsMock: BillRecord[] = [
  {
    id: "CB-2026-0418",
    customerName: "Northwind Logistics",
    billingMonth: "2026-04",
    amount: 126500,
    status: "ready",
    channel: "EDI",
    updatedAt: "2026-04-24 09:18",
  },
  {
    id: "CB-2026-0415",
    customerName: "Helio Retail Group",
    billingMonth: "2026-04",
    amount: 87420.35,
    status: "reviewing",
    channel: "Portal",
    updatedAt: "2026-04-24 08:46",
  },
  {
    id: "CB-2026-0409",
    customerName: "Blue Peak Energy",
    billingMonth: "2026-04",
    amount: 192880.9,
    status: "overdue",
    channel: "Email",
    updatedAt: "2026-04-23 18:05",
  },
  {
    id: "CB-2026-0402",
    customerName: "Nova Health Labs",
    billingMonth: "2026-03",
    amount: 45670,
    status: "paid",
    channel: "Portal",
    updatedAt: "2026-04-22 16:12",
  },
];

export const customersMock: CustomerSnapshot[] = [
  {
    id: "CUS-101",
    name: "Northwind Logistics",
    creditLevel: "A",
    activeContracts: 12,
    outstandingAmount: 126500,
    billingHealth: 91,
    primaryContact: "Olivia Chen",
  },
  {
    id: "CUS-204",
    name: "Helio Retail Group",
    creditLevel: "B+",
    activeContracts: 8,
    outstandingAmount: 87420.35,
    billingHealth: 76,
    primaryContact: "Mason Wu",
  },
  {
    id: "CUS-318",
    name: "Blue Peak Energy",
    creditLevel: "A-",
    activeContracts: 5,
    outstandingAmount: 192880.9,
    billingHealth: 58,
    primaryContact: "Sophia Lin",
  },
];

export const reconciliationTasksMock: ReconciliationTask[] = [
  {
    id: "REC-2401",
    bank: "Industrial Bank",
    period: "2026-04-24",
    matched: 182,
    total: 197,
    status: "running",
    owner: "Finance Ops",
  },
  {
    id: "REC-2398",
    bank: "ICBC",
    period: "2026-04-23",
    matched: 210,
    total: 210,
    status: "done",
    owner: "Finance Ops",
  },
  {
    id: "REC-2396",
    bank: "Bank of China",
    period: "2026-04-23",
    matched: 96,
    total: 122,
    status: "attention",
    owner: "Settlement Team",
  },
];
