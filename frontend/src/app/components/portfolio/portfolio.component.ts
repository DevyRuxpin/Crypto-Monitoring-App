import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { PortfolioService } from '../../core/services/portfolio.service';

@Component({
  selector: 'app-portfolio',
  templateUrl: './portfolio.component.html',
  styleUrls: ['./portfolio.component.scss']
})
export class PortfolioComponent implements OnInit {
  portfolioForm: FormGroup;
  portfolio: any[] = [];
  loading = false;

  constructor(
    private fb: FormBuilder,
    private portfolioService: PortfolioService,
    private snackBar: MatSnackBar
  ) {
    this.portfolioForm = this.fb.group({
      symbol: ['', Validators.required],
      quantity: ['', [Validators.required, Validators.min(0)]]
    });
  }

  ngOnInit(): void {
    this.loadPortfolio();
  }

  onSubmit(): void {
    if (this.portfolioForm.valid) {
      this.loading = true;
      const { symbol, quantity } = this.portfolioForm.value;
      
      this.portfolioService.addToPortfolio(symbol, quantity)
        .subscribe(
          () => {
            this.snackBar.open('Successfully added to portfolio', 'Close', { duration: 3000 });
            this.loadPortfolio();
            this.portfolioForm.reset();
          },
          error => {
            this.snackBar.open(error.message || 'Failed to add to portfolio', 'Close', { duration: 3000 });
          }
        )
        .add(() => this.loading = false);
    }
  }

  private loadPortfolio(): void {
    this.loading = true;
    this.portfolioService.getPortfolio()
      .subscribe(
        data => this.portfolio = data,
        error => this.snackBar.open(error.message || 'Failed to load portfolio', 'Close', { duration: 3000 })
      )
      .add(() => this.loading = false);
  }
}