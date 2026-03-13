export type Task = {
  id: string;
  date: string;
  status: string;
  created_at: string;
};

export type ProgressBar = {
  task_id: string;
  status: string;
  total_matches: number;
  current_match: number;
};
