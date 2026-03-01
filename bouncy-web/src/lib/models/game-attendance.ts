export interface GameAttendance {
  userId: string;
  userName?: string;
  checkedIn: boolean;
  status: number;
  checkInComment: string;
  createdAt: Date;
  updatedAt: Date;
}
