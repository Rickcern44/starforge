export interface GameAttendance {
  userId: string;
  checkedIn: boolean;
  status: number;
  checkInComment: string;
  createdAt: Date;
  updatedAt: Date;
}
