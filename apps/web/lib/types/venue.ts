export interface Location {
  latitude: number;
  longitude: number;
}

export interface Address {
  street: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
}

export interface Venue {
  id: string;
  name: string;
  location: Location;
  address: Address;
  venue_types: string[];
  description?: string;
  capacity?: number;
  rating?: number;
  external_source?: string;
  external_source_id?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateVenueRequest {
  name: string;
  location: Location;
  address: Address;
  venue_types: string[];
  description?: string;
  capacity?: number;
  external_source?: string;
  external_source_id?: string;
}

export interface UpdateVenueRequest {
  name?: string;
  location?: Location;
  address?: Address;
  venue_types?: string[];
  description?: string;
  capacity?: number;
}

export interface VenueSearchParams {
  venue_types?: string[];
  latitude?: number;
  longitude?: number;
  radius_km?: number;
  city?: string;
  state?: string;
  limit?: number;
  offset?: number;
}

export interface Booking {
  id: string;
  venue_id: string;
  user_id: string;
  event_name: string;
  start_time: string;
  end_time: string;
  status: 'pending' | 'confirmed' | 'cancelled';
  created_at: string;
  updated_at: string;
}

export interface CreateBookingRequest {
  venue_id: string;
  user_id: string;
  event_name: string;
  start_time: string;
  end_time: string;
}