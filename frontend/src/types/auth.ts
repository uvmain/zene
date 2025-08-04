export interface SessionCheck {
  loggedIn: boolean
  isAdmin: boolean
}

export interface User {
  id: number
  username: string
  is_admin: boolean
  created_at: string
}

export interface UsersResponse {
  users: User[]
  status: string
}

export interface TokenResponse {
  token: string
  expires_at: string
}
