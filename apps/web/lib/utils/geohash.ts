/**
 * Geohash utility functions for venue search
 */

const BASE32 = '0123456789bcdefghjkmnpqrstuvwxyz';

export function encodeGeohash(latitude: number, longitude: number, precision: number = 5): string {
  let lat = latitude;
  let lng = longitude;
  let latRange = [-90.0, 90.0];
  let lngRange = [-180.0, 180.0];
  
  let geohash = '';
  let bits = 0;
  let bit = 0;
  let even = true;
  
  while (geohash.length < precision) {
    if (even) {
      // longitude
      const mid = (lngRange[0] + lngRange[1]) / 2;
      if (lng >= mid) {
        bit = (bit << 1) + 1;
        lngRange[0] = mid;
      } else {
        bit = bit << 1;
        lngRange[1] = mid;
      }
    } else {
      // latitude
      const mid = (latRange[0] + latRange[1]) / 2;
      if (lat >= mid) {
        bit = (bit << 1) + 1;
        latRange[0] = mid;
      } else {
        bit = bit << 1;
        latRange[1] = mid;
      }
    }
    
    even = !even;
    
    if (++bits === 5) {
      geohash += BASE32[bit];
      bits = 0;
      bit = 0;
    }
  }
  
  return geohash;
}

export function getGeohashPrefixes(latitude: number, longitude: number, radiusKm: number): string[] {
  // For simplicity, we'll use different precision levels based on radius
  let precision: number;
  if (radiusKm <= 1) precision = 7;
  else if (radiusKm <= 5) precision = 6;
  else if (radiusKm <= 25) precision = 5;
  else if (radiusKm <= 125) precision = 4;
  else precision = 3;
  
  const centerHash = encodeGeohash(latitude, longitude, precision);
  
  // For a more complete implementation, we would calculate neighboring geohashes
  // For now, we'll return the center hash and some variations
  const prefixes = [centerHash];
  
  // Add some neighboring prefixes by modifying the last character
  if (precision > 1) {
    const baseHash = centerHash.slice(0, -1);
    for (let i = 0; i < BASE32.length; i++) {
      const neighborHash = baseHash + BASE32[i];
      if (neighborHash !== centerHash) {
        prefixes.push(neighborHash);
      }
    }
  }
  
  return prefixes;
}