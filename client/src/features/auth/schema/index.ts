import z from "zod"


export const registerSchema = z.object({
  email: z.email('Please enter a valid email address').min(1, 'Email is required'),
  password: z.string().min(8, 'Password must be at least 8 characters long').min(1, 'Password is required'),
})

export const loginSchema = z.object({
  email: z.email('Please enter a valid email address').min(1, 'Email is required'),
  password: z.string().min(8, 'Password must be at least 8 characters long').min(1, 'Password is required'),
})

export type RegisterFormData = z.infer<typeof registerSchema>
export type LoginFormData = z.infer<typeof loginSchema>