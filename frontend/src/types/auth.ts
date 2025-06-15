export interface SessionCheck {
  loggedIn: boolean
}

export interface User {
  id: number
  username: string
  is_admin: boolean
  created_at: string
}
