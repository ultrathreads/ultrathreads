// src/types/admin.ts

export interface SystemInfo {
  appName: string;
  appVersion: string;
  arch: string;
  buildCommit: string;
  buildTime: string;
  goversion: string;
  numCpu: number;
  os: string;
  upTime: string;
}