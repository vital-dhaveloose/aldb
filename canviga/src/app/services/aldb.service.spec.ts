import { TestBed } from '@angular/core/testing';

import { AldbService } from './aldb.service';

describe('AldbService', () => {
  let service: AldbService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AldbService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
